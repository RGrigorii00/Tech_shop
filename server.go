package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Role      string  `json:"role"`
	Gender    string  `json:"gender"`
	Address   string  `json:"address"`
	Photo_url string  `json:"photo_url"`
	Balance   float32 `json:"balance"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	Picture_url string  `json:picture_url`
}

var db *sql.DB

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func initDB() (*sql.DB, error) {
	psqlInfo := "host=localhost port=5432 user=postgres password=superuser dbname=Products sslmode=disable"
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func downloadFile(c echo.Context) error {
	filePath := "./NAZACHET.rar"
	fileName := "your-archive.rar"

	file, err := os.Open(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при открытии файла")
	}
	defer file.Close()

	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName)
	c.Response().Header().Set(echo.HeaderContentType, "application/x-rar-compressed")

	return c.Stream(http.StatusOK, "application/x-rar-compressed", file)
}

func getPaginationParams(c echo.Context) (int, int, int) {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page := 1
	limit := 9
	var err error

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			page = 1
		}
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			limit = 2
		}
	}

	offset := (page - 1) * limit
	return page, limit, offset
}

func getProducts(c echo.Context) error {
	page, limit, offset := getPaginationParams(c)

	var totalRecords int
	err := db.QueryRow("SELECT COUNT(*) FROM public.\"Products\"").Scan(&totalRecords)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count records"})
	}

	// Рассчитываем общее количество страниц
	totalPages := (totalRecords + limit - 1) / limit // Округление вверх

	// Проверяем, чтобы страница не выходила за пределы
	if page > totalPages {
		page = totalPages
	}
	if page < 1 {
		page = 1
	}

	offset = (page - 1) * limit

	rows, err := db.Query("SELECT id, name, description, price, stock, category, picture_url FROM public.\"Products\" ORDER BY id ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query the database"})
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.Picture_url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan row"})
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while iterating rows"})
	}

	response := map[string]interface{}{
		"products":    products,
		"total_pages": totalPages,
		"page":        page,
	}

	return c.JSON(http.StatusOK, response)
}

func getProductstable(c echo.Context) error {
	page, limit, offset := getPaginationParams(c)
	log.Printf("Current page: %d", page)

	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM public.\"Products\"").Scan(&totalCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get total count of products"})
	}

	// Вычисляем общее количество страниц
	totalPages := (totalCount + limit - 1) / limit

	// Ограничиваем текущую страницу на максимальное количество страниц
	if page > totalPages {
		page = totalPages
	}

	rows, err := db.Query("SELECT id, name, description, price, stock, category, picture_url FROM public.\"Products\" ORDER BY id ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query the database"})
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.Picture_url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan row"})
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while iterating rows"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"products":    products,
		"page":        page,
		"total_pages": totalPages,
	})
}

type ProductID struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	Created_at  time.Time `json:"created_at"`
	Picture_url string    `json:"picture_url"`
	Model_url   string    `json:"model_url"`
	Texture_url string    `json:"texture_url"`
}

func getProductByID(c echo.Context) error {
	id := c.Param("id")

	var product ProductID
	err := db.QueryRow("SELECT id, name, description, price, stock, category, created_at, picture_url, COALESCE(model_url, '../../3dmodels/Chair.glb') AS model_url, COALESCE(texture_url, '../../3dmodels/normalsidyshkatexture_png.png') AS texture_url FROM public.\"Products\" WHERE id = $1;", id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.Category,
		&product.Created_at,
		&product.Picture_url,
		&product.Model_url,
		&product.Texture_url,
	)
	if err != nil {
		log.Printf("Failed to get product by ID: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}

func createSessionInDB(userID int) (string, error) {
	sessionID := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())

	// Сохраняем сессию в базе данных
	_, err := db.Exec(`INSERT INTO sessions1 (session_id, user_id) VALUES ($1, $2)`, sessionID, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func getUserBySession(sessionID string) (int, error) {
	var userID int
	err := db.QueryRow(`SELECT user_id FROM sessions1 WHERE session_id = $1`, sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func register(c echo.Context) error {
	var input struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Role      string  `json:"role"`
		Gender    string  `json:"gender"`
		Address   string  `json:"address"`
		Photo_url string  `json:"photo_url"`
		Balance   float64 `json:"balance"`
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат данных")
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Все поля обязательны для заполнения")
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users1 WHERE email=$1)", input.Email).Scan(&exists)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при проверке пользователя")
	}
	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "Пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при хешировании пароля")
	}

	_, err = db.Exec(`
		INSERT INTO users1 (
			username, email, password, role, gender, address, photo_url, balance
		) 
		VALUES (
		    CAST($1 AS CHAR(255)), 
    $2, 
    CAST($3 AS TEXT), 
    'user', 
    'M', 
    'Address', 
    'https://via.placeholder.com/150', 
    1000000000
		);
	`,
		input.Username,
		input.Email,
		string(hashedPassword),
	)

	if err != nil {
		fmt.Println("Ошибка при добавлении пользователя:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при добавлении пользователя")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	})
}

func login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат данных")
	}

	if input.Email == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email и пароль обязательны")
	}

	var userID int
	var dbPassword string
	err := db.QueryRow("SELECT id, password FROM users1 WHERE email=$1", input.Email).Scan(&userID, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Неверный email или пароль")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при проверке данных пользователя")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(input.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Неверный email или пароль")
	}

	var userName, userAvatar string
	err = db.QueryRow("SELECT id, username, photo_url FROM users1 WHERE id=$1", userID).Scan(&userID, &userName, &userAvatar)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных пользователя")
	}

	// Создаем сессию
	sessionID, err := createSessionInDB(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при создании сессии")
	}

	// Сохраняем сессию в куки
	session, _ := store.Get(c.Request(), "session")
	session.Values["sessionID"] = sessionID
	session.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Авторизация успешна",
		"data": map[string]interface{}{
			"id":        userID,
			"username":  userName,
			"email":     input.Email,
			"photo_url": userAvatar,
			"session":   sessionID,
		},
	})
}

func profile(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	sessionID, ok := session.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Необходимо авторизоваться")
	}

	var sessionData struct {
		Id        string `json:"id"`
		SessionID string `json:"session_id"`
		UserID    int    `json:"user_id"`
		CreatedAt string `json:"created_at"`
	}
	err := db.QueryRow(`SELECT id, session_id, user_id, created_at FROM sessions1 WHERE session_id=$1`, sessionID).
		Scan(&sessionData.Id, &sessionData.SessionID, &sessionData.UserID, &sessionData.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Сессия недействительна")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных сессии")
	}

	var userData struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Role      string  `json:"role"`
		Gender    string  `json:"gender"`
		Address   string  `json:"address"`
		Photo_url string  `json:"photo_url"`
		Balance   float32 `json:"balance"`
	}
	err = db.QueryRow(`SELECT id, username, email, password , role, gender, address, photo_url , balance FROM users1 WHERE id=$1`, sessionData.UserID).
		Scan(&userData.ID, &userData.Username, &userData.Email, &userData.Password, &userData.Role, &userData.Gender, &userData.Address, &userData.Photo_url, &userData.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Пользователь не найден")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных пользователя")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":    userData,
		"session": sessionData,
	})
}

func profile15(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	sessionID, ok := session.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Необходимо авторизоваться")
	}

	var sessionData struct {
		SessionID string `json:"session_id"`
		UserID    int    `json:"user_id"`
	}
	err := db.QueryRow(`SELECT session_id, user_id FROM sessions1 WHERE session_id=$1`, sessionID).
		Scan(&sessionData.SessionID, &sessionData.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Сессия недействительна")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных сессии")
	}

	var userData struct {
		Username string  `json:"username"`
		PhotoURL string  `json:"photo_url"`
		Balance  float32 `json:"balance"`
	}
	err = db.QueryRow(`SELECT username, photo_url, balance FROM users1 WHERE id=$1`, sessionData.UserID).
		Scan(&userData.Username, &userData.PhotoURL, &userData.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Пользователь не найден")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных пользователя")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": userData,
	})
}

func logout(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	sessionID, ok := session.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Необходимо авторизоваться")
	}

	// Удаляем сессию из базы данных
	_, err := db.Exec("DELETE FROM sessions1 WHERE session_id=$1", sessionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при выходе из системы")
	}

	// Удаляем сессию из браузера (удаляем куки)
	session.Options = &sessions.Options{
		MaxAge: -1, // Устанавливаем отрицательное время жизни куки, чтобы она была удалена
	}
	session.Values["sessionID"] = nil       // Очистить значение сессии
	session.Save(c.Request(), c.Response()) // Сохранить изменения

	return c.String(http.StatusOK, "Выход успешен")
}

func uploadPhoto(c echo.Context) error {
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Ошибка при получении файла",
		})
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при открытии файла",
		})
	}
	defer src.Close()

	// Получаем ID пользователя
	userID := c.FormValue("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID пользователя не указан",
		})
	}

	// Генерируем уникальное имя файла
	originalFileName := filepath.Base(file.Filename)
	extension := filepath.Ext(originalFileName)
	fileName := fmt.Sprintf("%s_%d%s", userID, time.Now().UnixNano(), extension)
	filePath := filepath.Join("src/assets/users_photos/", fileName)

	// Сохраняем файл
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при сохранении файла",
		})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при копировании файла",
		})
	}

	// Формируем URL
	photoURL := fmt.Sprintf("src/assets/users_photos/%s", fileName)

	query := "UPDATE users1 SET photo_url = $1 WHERE id = $2"

	_, err = db.Exec(query, photoURL, userID)
	if err != nil {
		log.Printf("Database error: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при обновлении базы данных",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"photo_url": photoURL,
	})
}

type CartItem struct {
	CartID       int     `json:"cart_id"`
	UserName     string  `json:"user_name"`
	ProductName  string  `json:"product_name"`
	ProductId    int     `json:"product_id"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
}

func getCartItems(c echo.Context) error {
	userID := c.QueryParam("user_id")

	var cartItems []CartItem

	rows, err := db.Query("SELECT ci.id as catr_id, u.username AS user_name, p.name AS product_name, ci.product_id, p.price AS product_price, ci.quantity FROM cart_items ci JOIN users1 u ON ci.user_id = u.id JOIN public.\"Products\" p ON ci.product_id = p.id WHERE ci.user_id = $1", userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при загрузке корзины"})
	}
	defer rows.Close()

	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.CartID, &item.UserName, &item.ProductName, &item.ProductId, &item.ProductPrice, &item.Quantity); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при чтении данных"})
		}
		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при обработке результата запроса"})
	}

	return c.JSON(http.StatusOK, cartItems)
}

func addToCart(c echo.Context) error {
	var input struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректный запрос")
	}

	_, err := db.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity",
		input.UserID, input.ProductID, input.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Сервер Ошибка при добавлении в корзину")
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Товар добавлен в корзину"})
}

func removeFromCart(c echo.Context) error {
	cartID := c.QueryParam("id")

	cartIDInt, err := strconv.Atoi(cartID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректный cart_id")
	}

	result, err := db.Exec("DELETE FROM cart_items WHERE id = $1", cartIDInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при удалении товаров из корзины")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении количества удаленных строк")
	}

	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Корзина пуста или не найдена")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Товары удалены из корзины"})
}

func clearCart(c echo.Context) error {
	userID := c.QueryParam("user_id")

	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Отсутствует user_id")
	}

	result, err := db.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при очистке корзины")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении количества удаленных строк")
	}

	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Корзина пуста или не найдена")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Корзина очищена"})
}

func getAllBuys(c echo.Context) error {
	userID := c.Param("id")
	log.Printf("Получен userID: %s", userID)

	rows, err := db.Query(`
		SELECT o.id, u.username AS user_name, o.user_id, p.name AS product_name, 
		       o.product_id, o.quantity, o.total_price, 
		       TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
		FROM orders o
		JOIN users1 u ON o.user_id = u.id
		JOIN public."Products" p ON o.product_id = p.id
		WHERE o.user_id = $1
	`, userID)

	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка выполнения запроса"})
	}
	defer rows.Close()

	var sells []Sell

	for rows.Next() {
		var sell Sell
		if err := rows.Scan(&sell.ID, &sell.UserName, &sell.UserID, &sell.ProductName, &sell.ProductID, &sell.Quantity, &sell.TotalPrice, &sell.CreatedAt); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
		}
		sells = append(sells, sell)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при переборе строк: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
	}

	// Проверка, если нет результатов
	if len(sells) == 0 {
		log.Printf("Нет заказов для userID: %s", userID)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Заказы не найдены"})
	}

	return c.JSON(http.StatusOK, sells)
}

func getAllOrders(c echo.Context) error {
	rows, err := db.Query("SELECT o.id, username AS user_name,  user_id, name AS product_name, o.product_id AS product_id, o.quantity, o.total_price,     TO_CHAR(o.created_at, 'YYYY-MM-DD     HH24:MI:SS') AS created_at FROM orders o JOIN users1 u ON o.user_id = u.id JOIN public.\"Products\" p ON o.product_id = p.id;")
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка выполнения запроса"})
	}
	defer rows.Close()

	var sells []Sell // Структура Sell, аналогичная Order

	for rows.Next() {
		var sell Sell
		if err := rows.Scan(&sell.ID, &sell.UserName, &sell.UserID, &sell.ProductName, &sell.ProductID, &sell.Quantity, &sell.TotalPrice, &sell.CreatedAt); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
		}
		sells = append(sells, sell)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при переборе строк: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
	}

	return c.JSON(http.StatusOK, sells)
}

func editProduct(c echo.Context) error {
	var product Product

	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error: Invalid product ID: %s", id)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}
	log.Printf("Editing product with ID: %d", productID)

	// Получаем данные из формы
	product.Name = c.FormValue("Name")
	product.Description = c.FormValue("Description")
	log.Printf("Product name: %s, Description: %s", product.Name, product.Description)

	price, err := strconv.ParseFloat(c.FormValue("Price"), 64)
	if err != nil {
		log.Printf("Error: Invalid price value: %s", c.FormValue("Price"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid price value"})
	}
	product.Price = price
	log.Printf("Price: %.2f", product.Price)

	stock, err := strconv.Atoi(c.FormValue("Stock"))
	if err != nil {
		log.Printf("Error: Invalid stock value: %s", c.FormValue("Stock"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid stock value"})
	}
	product.Stock = stock
	log.Printf("Stock: %d", product.Stock)

	product.Category = c.FormValue("Category")
	log.Printf("Category: %s", product.Category)

	// Обработка изображения
	file, err := c.FormFile("Picture_url")
	if err == nil {
		log.Printf("Processing file upload: %s", file.Filename)

		// Загружаем изображение, если оно передано
		dirPath := "./assets/products_photos"
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Printf("Error: Failed to create directory %s: %v", dirPath, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
		}

		filePath := filepath.Join(dirPath, file.Filename)
		src, err := file.Open()
		if err != nil {
			log.Printf("Error: Failed to open file %s: %v", file.Filename, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error: Failed to save file %s: %v", file.Filename, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("Error: Failed to copy file %s: %v", file.Filename, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to copy file"})
		}

		product.Picture_url = fmt.Sprintf("/assets/products_photos/%s", file.Filename)
		log.Printf("Image uploaded successfully: %s", product.Picture_url)
	} else {
		// Если файл не был передан, оставляем старую картинку
		log.Println("No image uploaded, keeping existing picture_url.")
		err = db.QueryRow("SELECT picture_url FROM public.\"Products\" WHERE id = $1", productID).Scan(&product.Picture_url)
		if err != nil {
			log.Printf("Error: Failed to fetch current picture_url for product ID %d: %v", productID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch current picture_url"})
		}
		log.Printf("Current picture_url: %s", product.Picture_url)
	}

	log.Printf("Updating product in the database with ID: %d", productID)
	_, err = db.Exec(
		"UPDATE public.\"Products\" SET name = $1, description = $2, price = $3, stock = $4, category = $5, picture_url = $6 WHERE id = $7",
		product.Name, product.Description, product.Price, product.Stock, product.Category, product.Picture_url, productID,
	)
	if err != nil {
		log.Printf("Error: Failed to update product in database for ID %d: %v", productID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}

	log.Printf("Product with ID %d updated successfully", productID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

func addProduct(c echo.Context) error {
	var product Product

	product.Name = c.FormValue("Name")
	product.Description = c.FormValue("Description")
	price, err := strconv.ParseFloat(c.FormValue("Price"), 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid price value"})
	}
	product.Price = price

	stock, err := strconv.Atoi(c.FormValue("Stock"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid stock value"})
	}
	product.Stock = stock
	product.Category = c.FormValue("Category")

	// Обработка изображения
	file, err := c.FormFile("Picture_url")
	if err != nil {
		// Если изображения нет, используем placeholder, который тоже почему-то не доступен ЫЫЫЫЫЫЫЫЫЫЫЫЫЫЫЫЫ
		product.Picture_url = "https://via.placeholder.com/150"
	} else {
		// Загружаем изображение
		dirPath := "./assets/products_photos"
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create directory"})
		}

		filePath := filepath.Join(dirPath, file.Filename)
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to copy file"})
		}

		product.Picture_url = fmt.Sprintf("/assets/products_photos/%s", file.Filename)
	}

	var maxID int
	err = db.QueryRow("SELECT MAX(id) FROM public.\"Products\"").Scan(&maxID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get max product id"})
	}
	product.ID = maxID + 1

	_, err = db.Exec(
		"INSERT INTO public.\"Products\" (id, name, description, price, stock, category, created_at, picture_url) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)",
		product.ID, product.Name, product.Description, product.Price, product.Stock, product.Category, product.Picture_url,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add product"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":    "Product added successfully",
		"product_id": product.ID,
		"product":    product,
	})
}

func deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	// Удаляем товар по ID
	_, err = db.Exec("DELETE FROM public.\"Products\" WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func uploadProductPhoto(c echo.Context) error {
	// Получаем файл из формы
	file, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Ошибка при получении файла",
		})
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при открытии файла",
		})
	}
	defer src.Close()

	// Получаем ID товара
	productID := c.FormValue("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID товара не указан",
		})
	}

	// Генерируем уникальное имя файла
	originalFileName := filepath.Base(file.Filename)
	extension := filepath.Ext(originalFileName)
	fileName := fmt.Sprintf("product_%s_%d%s", productID, time.Now().UnixNano(), extension)
	filePath := filepath.Join("src/assets/products_photos/", fileName)

	// Сохраняем файл
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при сохранении файла",
		})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при копировании файла",
		})
	}

	// Формируем URL для изображения
	photoURL := fmt.Sprintf("src/assets/products_photos/%s", fileName)

	// Обновляем путь к изображению в базе данных
	query := "UPDATE products SET picture_url = $1 WHERE id = $2"
	_, err = db.Exec(query, photoURL, productID)
	if err != nil {
		log.Printf("Database error: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Ошибка при обновлении базы данных",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"picture_url": photoURL,
	})
}

type Sell struct {
	ID          int     `json:"sell_id"`
	UserID      int     `json:"sell_u_id"`
	UserName    string  `json:"sell_un_id"`
	ProductID   int     `json:"sell_p_id"`
	ProductName string  `json:"sell_pn_id"`
	Quantity    int     `json:"sell_q_id"`
	TotalPrice  float64 `json:"sell_tp_id"`
	CreatedAt   string  `json:"autoCreateTime"`
	CartID      int     `json:"sell_cid_id"`
}

type CartItem15 struct {
	CartId      int     `json:"cart_id"`
	ProductName string  `json:"product_name"`
	ProductId   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type PurchaseRequest struct {
	UserID    int          `json:"user_id"`
	CartItems []CartItem15 `json:"cart_items"`
}

func makePurchase(c echo.Context) error {
	var req PurchaseRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("Ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
	}

	// Проверяем, что корзина не пуста
	if len(req.CartItems) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Корзина пуста"})
	}

	// Проверка баланса пользователя
	var userBalance float64
	err := db.QueryRow("SELECT balance FROM users1 WHERE id = $1", req.UserID).Scan(&userBalance)
	if err != nil {
		log.Printf("Ошибка при получении баланса пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось найти пользователя"})
	}

	// Рассчитываем общую сумму заказа
	var totalOrderPrice float64
	for _, item := range req.CartItems {
		totalOrderPrice += item.TotalPrice
	}

	// Проверка наличия достаточных средств на балансе
	if userBalance < totalOrderPrice {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Недостаточно средств на балансе"})
	}

	// Транзакция
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при начале транзакции"})
	}
	defer tx.Rollback()

	for _, item := range req.CartItems {

		_, err := tx.Exec(`INSERT INTO orders (user_id, product_id, quantity, total_price, created_at, cart_id) VALUES ($1, $2, $3, $4, NOW(), $5)`,
			req.UserID, item.ProductId, item.Quantity, item.TotalPrice, item.CartId)
		if err != nil {
			log.Printf("Ошибка при добавлении записи в таблицу orders: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении записи в таблицу orders"})
		}
		_, err = tx.Exec("UPDATE public.\"Products\" SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductId)
		if err != nil {
			log.Printf("Ошибка при обновлении количества товара: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении количества товара"})
		}

		_, err = tx.Exec(`UPDATE users1 SET balance = balance - $1 WHERE id = $2`, item.TotalPrice, req.UserID)
		if err != nil {
			log.Printf("Ошибка при обновлении баланса пользователя: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении баланса пользователя"})
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при завершении транзакции"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Покупка успешно завершена"})
}

type CartData struct {
	Cart []CartItem `json:"cart"`
}

func generatePDF(c echo.Context) error {
	var cartData CartData
	fmt.Println(cartData)

	err := json.NewDecoder(c.Request().Body).Decode(&cartData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode JSON")
	}
	fmt.Println(cartData)

	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "", 14)

	columnWidths := []float64{
		20,  // Cart ID
		40,  // User Name
		100, // Product Name
		50,  // Price
		30,  // Quantity
		40,  // Total Price
	}

	pdf.Cell(columnWidths[0], 10, "Cart ID")
	pdf.Cell(columnWidths[1], 10, "User Name")
	pdf.Cell(columnWidths[2], 10, "Product Name")
	pdf.Cell(columnWidths[3], 10, "Price [RUB]")
	pdf.Cell(columnWidths[4], 10, "Quantity")
	pdf.Cell(columnWidths[5], 10, "Total Price [RUB]")
	pdf.Ln(10)

	for _, item := range cartData.Cart {
		pdf.Cell(columnWidths[0], 10, fmt.Sprintf("%d", item.CartID))
		pdf.Cell(columnWidths[1], 10, item.UserName)
		pdf.Cell(columnWidths[2], 10, item.ProductName)
		pdf.Cell(columnWidths[3], 10, fmt.Sprintf("%.2f ₽", item.ProductPrice))
		pdf.Cell(columnWidths[4], 10, fmt.Sprintf("%d", item.Quantity))
		pdf.Cell(columnWidths[5], 10, fmt.Sprintf("%.2f ₽", item.ProductPrice*float64(item.Quantity)))
		pdf.Ln(10)
	}

	// Отправляем PDF клиенту
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=sklad.pdf")
	err = pdf.Output(c.Response())
	if err != nil {
		log.Println("Error generating PDF:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating PDF")
	}

	return nil
}

type CartItem1 struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
}

type CartData1 struct {
	Cart []CartItem1 `json:"cart"`
}

func generatePDF15(c echo.Context) error {
	var cartData CartData1

	err := json.NewDecoder(c.Request().Body).Decode(&cartData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode JSON")
	}

	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "", 14)

	columnWidths := []float64{
		40,  // Product ID
		100, // Product Name
		50,  // Price [RUB]
		30,  // Quantity
		40,  // Total Price [RUB]
	}

	// Заголовки таблицы
	pdf.Cell(columnWidths[0], 10, "Product ID")
	pdf.Cell(columnWidths[1], 10, "Product Name")
	pdf.Cell(columnWidths[2], 10, "Price [RUB]")
	pdf.Cell(columnWidths[3], 10, "Quantity")
	pdf.Cell(columnWidths[4], 10, "Total Price [RUB]")
	pdf.Ln(10)

	for _, item := range cartData.Cart {
		// Для каждого товара выводим данные
		pdf.Cell(columnWidths[0], 10, fmt.Sprintf("%d", item.ID))
		pdf.Cell(columnWidths[1], 10, item.Name)
		pdf.Cell(columnWidths[2], 10, fmt.Sprintf("%.2f ₽", item.Price))
		pdf.Cell(columnWidths[3], 10, fmt.Sprintf("%d", item.Stock))
		pdf.Cell(columnWidths[4], 10, fmt.Sprintf("%.2f ₽", item.Price*float64(item.Stock)))
		pdf.Ln(10)
	}

	// Отправляем PDF клиенту
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=cart.pdf")
	err = pdf.Output(c.Response())
	if err != nil {
		log.Println("Error generating PDF:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating PDF")
	}

	return nil
}

func main() {
	_, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	e.GET("/daf", downloadFile) // Чтобы скачать с сервера

	e.GET("/buy/:id", getAllBuys) // Мои покупки

	e.POST("/generate-pdf", generatePDF)
	e.POST("/generate-pdf15", generatePDF15)

	e.Static("/src/assets/products_photos", "src/assets/products_photos")

	e.GET("/sells", getAllOrders) // Все покупки

	e.POST("/product_edit/:id", editProduct)        // Добавить новый товар
	e.POST("/products_add", addProduct)             // Добавить новый товар
	e.DELETE("/products_delete/:id", deleteProduct) // Удалить товар по ID

	e.POST("/make_purchase", makePurchase)
	e.GET("/get_cart", getCartItems)
	e.POST("/add_to_cart", addToCart)
	e.DELETE("/remove_from_cart", removeFromCart) // Удалить товар из корзины
	e.DELETE("/clear_cart", clearCart)

	e.GET("/products_table", getProductstable)
	e.GET("/products", getProducts)
	e.GET("/product/:id", getProductByID)

	e.POST("/register", register)
	e.POST("/login", login)
	e.GET("/profile15", profile15)
	e.GET("/profile", profile)
	e.GET("/logout", logout)

	e.POST("/upload-image", uploadProductPhoto)
	e.POST("/upload-photo", uploadPhoto)

	e.Logger.Fatal(e.Start(":8080"))
}
