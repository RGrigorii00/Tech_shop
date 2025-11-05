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

// Структура для товара
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	Picture_url string  `json:picture_url`
}

// Структура для соединения с базой данных
var db *sql.DB
var dbu *sql.DB

// Переменная для работы сессий
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Инициализация базы данных
func initDB() (*sql.DB, error) {
	psqlInfo := "host=79.174.88.80 port=16680 user=superuser password=Superuser_1 dbname=Products sslmode=disable"
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

// Инициализация базы данных
func initDBU() (*sql.DB, error) {
	psqlInfo := "host=79.174.88.80 port=16680 user=superuser password=Superuser_1 dbname=Products sslmode=disable"
	var err error
	dbu, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = dbu.Ping()
	if err != nil {
		return nil, err
	}

	return dbu, nil
}

func downloadFile(c echo.Context) error {
	filePath := "./NAZACHET.rar"   // Путь к вашему .rar файлу
	fileName := "your-archive.rar" // Имя файла, которое будет предложено при скачивании

	// Открыть .rar файл
	file, err := os.Open(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при открытии файла")
	}
	defer file.Close()

	// Устанавливаем заголовки для скачивания .rar файла
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName)
	c.Response().Header().Set(echo.HeaderContentType, "application/x-rar-compressed") // Указываем MIME-тип для .rar

	// Возвращаем файл
	return c.Stream(http.StatusOK, "application/x-rar-compressed", file)
}

// Функция для обработки пагинации
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

// Получить все товары с пагинацией
func getProducts(c echo.Context) error {
	page, limit, offset := getPaginationParams(c)
	log.Printf("Current page: %d", page)

	rows, err := db.Query("SELECT id, name, description, price, stock, category, picture_url FROM public.\"Products\" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		// log.Printf("Failed to query the database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query the database"})
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.Picture_url)
		if err != nil {
			// log.Printf("Failed to scan row: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan row"})
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		// log.Printf("Error while iterating rows: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while iterating rows"})
	}

	// Возвращаем ответ с товарами, количеством страниц и offset
	return c.JSON(http.StatusOK, products)
}

// Получить все товары с пагинацией
func getProductstable(c echo.Context) error {
	page, limit, offset := getPaginationParams(c)
	log.Printf("Current page: %d", page)

	rows, err := db.Query("SELECT id, name, description, price, stock, category, picture_url FROM public.\"Products\" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		// log.Printf("Failed to query the database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query the database"})
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.Picture_url)
		if err != nil {
			// log.Printf("Failed to scan row: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan row"})
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		// log.Printf("Error while iterating rows: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error while iterating rows"})
	}

	return c.JSON(http.StatusOK, products)
}

// Структура для товара
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

// Получить товар по ID
func getProductByID(c echo.Context) error {
	// Получаем ID товара из параметров URL
	id := c.Param("id")

	// Запрос для получения товара по ID
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

	fmt.Println(product)

	// Возвращаем товар в формате JSON
	return c.JSON(http.StatusOK, product)
}

// Создание сессии в базе данных
func createSessionInDB(userID int) (string, error) {
	sessionID := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())

	// Сохраняем сессию в базе данных
	_, err := db.Exec(`INSERT INTO sessions1 (session_id, user_id) VALUES ($1, $2)`, sessionID, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Получение пользователя по session ID из базы данных
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
		Balance   float64 `json:"balance"` // Используем float64 для точности
	}

	// Парсинг JSON
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат данных")
	}

	// Проверяем, что поля не пустые
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Все поля обязательны для заполнения")
	}

	// Проверяем, что пользователь уже существует
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

	// Добавляем нового пользователя в базу данных
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

// Авторизация пользователя
func login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Парсинг JSON
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат данных")
	}

	// Проверяем, что поля не пустые
	if input.Email == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email и пароль обязательны")
	}

	// Проверяем данные пользователя
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

	// Дополнительный запрос для получения всех данных пользователя
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

	// Отправляем ответ с данными пользователя
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Авторизация успешна",
		"data": map[string]interface{}{
			"id":        userID,
			"username":  userName,
			"email":     input.Email, // или userEmail, если оно было загружено
			"photo_url": userAvatar,
			"session":   sessionID, // отправляем sessionID для дальнейших запросов
		},
	})
}

// Получение профиля пользователя
func profile(c echo.Context) error {
	// Получаем сессию
	session, _ := store.Get(c.Request(), "session")
	sessionID, ok := session.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Необходимо авторизоваться")
	}

	// Получаем данные о сессии из базы данных
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

	// Получаем данные пользователя из базы данных
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
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Пользователь не найден")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении данных пользователя")
	}

	// Возвращаем данные пользователя и сессии
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":    userData,
		"session": sessionData,
	})
}

// Получение профиля пользователя
func profile15(c echo.Context) error {
	// Получаем сессию
	session, _ := store.Get(c.Request(), "session")
	sessionID, ok := session.Values["sessionID"].(string)
	if !ok || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Необходимо авторизоваться")
	}

	// Получаем данные о сессии из базы данных
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

	// Получаем только необходимые данные пользователя
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

	// Возвращаем только необходимые данные: аватарку, имя и баланс
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": userData,
	})
}

// Выход из системы
func logout(c echo.Context) error {
	// Получаем сессию из браузера
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
	// Получаем файл из формы
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

	// Обновляем базу данных
	query := "UPDATE users1 SET photo_url = $1 WHERE id = $2"
	// log.Printf("Executing query: %s with values photo_url: %s, user_id: %s", query, photoURL, userID)

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
	userID := c.QueryParam("user_id") // Получаем user_id из запроса
	// fmt.Println(userID)

	var cartItems []CartItem

	// Выполняем запрос с объединением таблиц
	rows, err := db.Query("SELECT ci.id as catr_id, u.username AS user_name, p.name AS product_name, ci.product_id, p.price AS product_price, ci.quantity FROM cart_items ci JOIN users1 u ON ci.user_id = u.id JOIN public.\"Products\" p ON ci.product_id = p.id WHERE ci.user_id = $1", userID)

	// fmt.Println(rows)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при загрузке корзины"})
	}
	defer rows.Close()

	// Читаем данные из результата запроса
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

	// fmt.Println(cartItems)
	return c.JSON(http.StatusOK, cartItems)
}

// Добавление товара в корзину
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

// Удаление товара из корзины по cart_id
func removeFromCart(c echo.Context) error {
	// Получаем cart_id из строки запроса
	cartID := c.QueryParam("id") // Получаем cart_id из строки запроса

	// Преобразуем cartID в тип int
	cartIDInt, err := strconv.Atoi(cartID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Некорректный cart_id")
	}

	// Выполняем запрос на удаление всех товаров из корзины по cart_id
	result, err := db.Exec("DELETE FROM cart_items WHERE id = $1", cartIDInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при удалении товаров из корзины")
	}

	// Проверяем, был ли удален хотя бы один товар
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении количества удаленных строк")
	}

	// Если товары не были удалены, выводим сообщение
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Корзина пуста или не найдена")
	}

	// Возвращаем успешный ответ
	return c.JSON(http.StatusOK, echo.Map{"message": "Товары удалены из корзины"})
}

// Очистка корзины
func clearCart(c echo.Context) error {
	// Получаем user_id из строки запроса
	userID := c.QueryParam("user_id")
	fmt.Println(userID)

	// Проверяем, что user_id не пустой
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Отсутствует user_id")
	}

	// Выполняем запрос на удаление всех товаров из корзины для указанного user_id
	result, err := db.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при очистке корзины")
	}

	// Получаем количество затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при получении количества удаленных строк")
	}

	// Если ничего не было удалено, выводим сообщение
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Корзина пуста или не найдена")
	}

	// Возвращаем успешный ответ
	return c.JSON(http.StatusOK, echo.Map{"message": "Корзина очищена"})
}

func getAllBuys(c echo.Context) error {
	// Получаем userID из параметров запроса
	userID := c.Param("id")
	fmt.Println(userID)
	// Логируем userID для диагностики
	log.Printf("Получен userID: %s", userID)

	// Выполняем запрос для получения заказов по userID
	rows, err := db.Query(`
		SELECT o.id, u.username AS user_name, o.user_id, p.name AS product_name, 
		       o.product_id, o.quantity, o.total_price, 
		       TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
		FROM orders o
		JOIN users1 u ON o.user_id = u.id
		JOIN public."Products" p ON o.product_id = p.id
		WHERE o.user_id = $1
	`, userID)

	// Проверка на ошибку запроса
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка выполнения запроса"})
	}
	defer rows.Close()

	// Массив для хранения заказов
	var sells []Sell

	// Перебираем строки результата запроса
	for rows.Next() {
		var sell Sell
		// Сканируем данные из строки в структуру
		if err := rows.Scan(&sell.ID, &sell.UserName, &sell.UserID, &sell.ProductName, &sell.ProductID, &sell.Quantity, &sell.TotalPrice, &sell.CreatedAt); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
		}
		// Добавляем заказ в массив
		sells = append(sells, sell)
	}

	// Проверка на наличие ошибок после перебора строк
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при переборе строк: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обработке данных"})
	}

	// Проверка, если нет результатов
	if len(sells) == 0 {
		log.Printf("Нет заказов для userID: %s", userID)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Заказы не найдены"})
	}

	// Возвращаем заказанные данные в формате JSON
	return c.JSON(http.StatusOK, sells)
}

func getAllOrders(c echo.Context) error {
	// Выполняем запрос для получения всех продаж
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

	// Получаем ID товара из параметров маршрута
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	// Получаем данные из формы
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
	if err == nil {
		// Загружаем изображение, если оно передано
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

	// Если файл не был передан, оставляем старую картинку
	if product.Picture_url == "" {
		err = db.QueryRow("SELECT picture_url FROM public.\"Products\" WHERE id = $1", productID).Scan(&product.Picture_url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch current picture_url"})
		}
	}

	// Обновляем товар в базе данных
	_, err = db.Exec(
		"UPDATE public.\"Products\" SET name = $1, description = $2, price = $3, stock = $4, category = $5, picture_url = $6 WHERE id = $7",
		product.Name, product.Description, product.Price, product.Stock, product.Category, product.Picture_url, productID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}

	// Возвращаем обновленный товар
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

func addProduct(c echo.Context) error {
	var product Product

	// Получаем данные из формы
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
		// Если изображения нет, используем placeholder
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

		product.Picture_url = fmt.Sprintf("../src/assets/products_photos/%s", file.Filename)
	}

	// Получение следующего ID
	var maxID int
	err = db.QueryRow("SELECT MAX(id) FROM public.\"Products\"").Scan(&maxID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get max product id"})
	}
	product.ID = maxID + 1

	// Вставляем товар в базу данных
	_, err = db.Exec(
		"INSERT INTO public.\"Products\" (id, name, description, price, stock, category, created_at, picture_url) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)",
		product.ID, product.Name, product.Description, product.Price, product.Stock, product.Category, product.Picture_url,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add product"})
	}

	// Возвращаем сообщение об успешном добавлении товара
	return c.JSON(http.StatusOK, echo.Map{
		"message":    "Product added successfully",
		"product_id": product.ID,
		"product":    product,
	})
}

func updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	product.ID = id

	// Запрос на обновление товара с возвратом актуальных данных
	err = db.QueryRow("UPDATE public.\"Products\" SET name = $1, description = $2, price = $3, stock = $4, category = $5, picture_url = $6 WHERE id = $7 RETURNING id, name, description, price, stock, category, picture_url", product.Name, product.Description, product.Price, product.Stock, product.Category, product.Picture_url, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category, &product.Picture_url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}

	// Возвращаем обновленный продукт
	return c.JSON(http.StatusOK, product)
}

// Удалить товар по ID
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

	// Возвращаем URL изображения
	return c.JSON(http.StatusOK, map[string]string{
		"picture_url": photoURL,
	})
}

type Order struct {
	ID         int `json:"primaryKey"`
	UserID     int
	ProductID  int
	Quantity   int
	TotalPrice float64
	CreatedAt  int64 `json:"autoCreateTime"`
	CartID     int
	CartItems  []CartItem
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

// Пример исправленного обработчика
func makePurchase(c echo.Context) error {
	var req PurchaseRequest

	// Привязка данных из тела запроса
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

	// Транзакция для атомарности операции
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

	// // Привязка JSON из тела запроса к структуре Product
	// if err := c.Bind(&itog); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	// }

	// // Теперь объект product содержит данные из тела запроса
	// fmt.Println("Полученный продукт:", itog)

	// Декодируем JSON-данные из запроса
	err := json.NewDecoder(c.Request().Body).Decode(&cartData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode JSON")
	}
	fmt.Println(cartData)

	// Создаем новый PDF документ с альбомной ориентацией
	pdf := gofpdf.New("L", "mm", "A4", "")

	// Добавляем страницу
	pdf.AddPage()

	// Устанавливаем шрифт Arial (стандартный шрифт с поддержкой кириллицы)
	pdf.SetFont("Arial", "", 14)

	// Распределим ширину на 6 столбцов
	columnWidths := []float64{
		20,  // Cart ID
		40,  // User Name
		100, // Product Name
		50,  // Price
		30,  // Quantity
		40,  // Total Price
	}

	// Расположим заголовки таблицы
	pdf.Cell(columnWidths[0], 10, "Cart ID")
	pdf.Cell(columnWidths[1], 10, "User Name")
	pdf.Cell(columnWidths[2], 10, "Product Name")
	pdf.Cell(columnWidths[3], 10, "Price [RUB]")
	pdf.Cell(columnWidths[4], 10, "Quantity")
	pdf.Cell(columnWidths[5], 10, "Total Price [RUB]")
	pdf.Ln(10)

	// Перебираем все элементы корзины и добавляем их в PDF
	for _, item := range cartData.Cart {
		pdf.Cell(columnWidths[0], 10, fmt.Sprintf("%d", item.CartID))
		pdf.Cell(columnWidths[1], 10, item.UserName)    // Просто используем строку без преобразования
		pdf.Cell(columnWidths[2], 10, item.ProductName) // Просто используем строку без преобразования
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

type SellData struct {
	Sells []Sell `json:"sells"` // Массив заказов
}

type ProductS struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Stock       string `json:"stock"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	Picture_url string `json:"picture_url"`
}

func generatePDFSell(c echo.Context) error {
	var sellData SellData

	// Декодирование JSON-данных из запроса
	err := json.NewDecoder(c.Request().Body).Decode(&sellData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode JSON")
	}

	// Логирование распакованных данных для проверки
	log.Printf("Полученные данные о заказах: %+v\n", sellData)

	// Проверка наличия заказов
	if len(sellData.Sells) == 0 {
		log.Println("Нет заказов для генерации PDF.")
		return echo.NewHTTPError(http.StatusBadRequest, "No orders to generate PDF")
	}

	// Создание нового PDF документа с альбомной ориентацией
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Устанавливаем шрифт Arial
	pdf.SetFont("Arial", "", 14)

	// Ширина колонок для таблицы
	columnWidths := []float64{
		20, // Sell ID
		40, // User ID
		60, // User Name
		40, // Product ID
		80, // Product Name
		30, // Quantity
		40, // Total Price [RUB]
		40, // Sale Date
		30, // Cart ID
	}

	// Заголовки таблицы
	headers := []string{
		"Sell ID", "User ID", "User Name", "Product ID",
		"Product Name", "Quantity", "Total Price [RUB]", "Sale Date", "Cart ID",
	}

	// Вывод заголовков таблицы
	for i, header := range headers {
		pdf.CellFormat(columnWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Перебираем все заказы и добавляем их в PDF
	for _, order := range sellData.Sells {
		pdf.CellFormat(columnWidths[0], 10, fmt.Sprintf("%d", order.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[1], 10, fmt.Sprintf("%d", order.UserID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[2], 10, order.UserName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(columnWidths[3], 10, fmt.Sprintf("%d", order.ProductID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[4], 10, order.ProductName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(columnWidths[5], 10, fmt.Sprintf("%d", order.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[6], 10, fmt.Sprintf("%.2f", order.TotalPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(columnWidths[7], 10, order.CreatedAt, "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[8], 10, fmt.Sprintf("%d", order.CartID), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Отправка PDF клиенту
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=orders.pdf")
	err = pdf.Output(c.Response())
	if err != nil {
		log.Println("Ошибка генерации PDF:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating PDF")
	}

	log.Println("PDF успешно сгенерирован и отправлен клиенту")
	return nil
}

func generatePDF15(c echo.Context) error {
	var cartData CartData1

	// Декодируем JSON-данные из запроса
	err := json.NewDecoder(c.Request().Body).Decode(&cartData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode JSON")
	}

	// Создаем новый PDF документ с альбомной ориентацией
	pdf := gofpdf.New("L", "mm", "A4", "")

	// Добавляем страницу
	pdf.AddPage()

	// Устанавливаем шрифт Arial (стандартный шрифт с поддержкой кириллицы)
	pdf.SetFont("Arial", "", 14)

	// Распределим ширину на 5 столбцов
	columnWidths := []float64{
		40,  // Product ID
		100, // Product Name
		50,  // Price [RUB]
		30,  // Quantity
		40,  // Total Price [RUB]
	}

	// Расположим заголовки таблицы
	pdf.Cell(columnWidths[0], 10, "Product ID")
	pdf.Cell(columnWidths[1], 10, "Product Name")
	pdf.Cell(columnWidths[2], 10, "Price [RUB]")
	pdf.Cell(columnWidths[3], 10, "Quantity")
	pdf.Cell(columnWidths[4], 10, "Total Price [RUB]")
	pdf.Ln(10)

	// Перебираем все элементы корзины и добавляем их в PDF
	for _, item := range cartData.Cart {
		// Для каждого товара выводим данные
		pdf.Cell(columnWidths[0], 10, fmt.Sprintf("%d", item.ID))
		pdf.Cell(columnWidths[1], 10, item.Name) // Просто используем строку без преобразования
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

func makePurchase15(c echo.Context) error {
	// Логируем тело запроса для диагностики
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Ошибка чтения тела запроса: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при чтении тела запроса"})
	}
	log.Printf("Тело запроса: %s", string(bodyBytes))

	var req PurchaseRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
	}

	log.Printf("Полученные данные: %+v", req)

	// Продолжение обработки...
	return c.JSON(http.StatusOK, map[string]string{"message": "Покупка успешно завершена"})
}

func main() {
	_, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	defer db.Close()

	dbu, err := initDBU()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	defer dbu.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Настраиваем CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Регистрируем обработчик маршрута для генерации PDF
	e.GET("/daf", downloadFile)

	// Регистрируем обработчик маршрута для генерации PDF
	// e.POST("/generate-pdf-product", generatePDFProduct)

	// Регистрируем обработчик маршрута для генерации PDF
	e.GET("/buy/:id", getAllBuys)

	// Регистрируем обработчик маршрута для генерации PDF
	e.POST("/generate-pdf-sell", generatePDFSell)

	// Регистрируем обработчик маршрута для генерации PDF
	e.POST("/generate-pdf", generatePDF)
	// Регистрируем обработчик маршрута для генерации PDF
	e.POST("/generate-pdf15", generatePDF15)

	// Этот обработчик должен быть настроен для пути "/make_purchase"
	e.POST("/make_purchase", makePurchase)

	// Разрешаем доступ к папке, где хранятся изображения
	e.Static("/src/assets/products_photos", "src/assets/products_photos")

	// Маршрут для загрузки изображений
	e.POST("/upload-image", uploadProductPhoto)

	// Роуты

	e.GET("/sells", getAllOrders) // Добавить новый товар

	e.POST("/product_edit/:id", editProduct) // Добавить новый товар
	// e.GET("/products_add", getProducts)      // Получить все товары с пагинацией
	e.POST("/products_add", addProduct)             // Добавить новый товар
	e.PUT("/products_update/:id", updateProduct)    // Обновить товар по ID
	e.DELETE("/products_delete/:id", deleteProduct) // Удалить товар по ID

	e.GET("/get_cart", getCartItems)
	e.POST("/add_to_cart", addToCart)

	e.DELETE("/remove_from_cart", removeFromCart) // Роут для удаления товара по cart_id

	e.DELETE("/clear_cart", clearCart)

	e.GET("/products_table", getProductstable)

	e.GET("/products", getProducts)
	e.GET("/product/:id", getProductByID)

	// Роуты
	e.POST("/register", register)
	e.POST("/login", login)

	e.GET("/profile15", profile15)
	e.GET("/profile", profile)
	e.GET("/logout", logout)

	// Маршрут для загрузки изображения
	e.POST("/upload-photo", uploadPhoto)

	e.Logger.Fatal(e.Start(":8080"))
}
