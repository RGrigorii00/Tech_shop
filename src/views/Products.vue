<template>

  <link href="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.css" rel="stylesheet">


  <h1 style="display: flex; justify-content: center; margin: 30px;">Товары</h1>
  <div>
    <!-- <router-link to="/products">
      <h5>-> Список товаров <-</h5>
    </router-link> -->
    <!-- <router-link to="/products_table">
      <h5>Таблица товаров</h5>
    </router-link> -->

    <div v-if="loading">Загрузка...</div>
    <div v-else style="margin-top: 30px;">

      <div class="row">
        <div class="col-12 col-sm-6 col-md-4 col-lg-3" v-for="product in products" :key="product.id">
          <div class="card card14" style="max-height: 500px; height: 500px;">
            <img :src="product.Picture_url" alt="product.name" class="card-img-top" />
            <div class="card-body d-flex flex-column">
              <h5 class="card-title">
                <router-link :to="'/product/' + product.id">{{ product.name }}</router-link>
              </h5>
              <p class="card-text text-truncate">{{ product.description }}</p>
              <p class="card-text">ID: {{ product.id }}</p>
              <p class="card-text">Цена: ₽ {{ product.price }}</p>
              <p class="card-text">Осталось на складе: {{ product.stock }}</p>


              <button @click="addToCart(product.id)" class="btn btn-primary mt-auto">Добавить в корзину</button>
              <!-- <button @click="editProduct(product)" class="btn btn-warning mt-2" v-if="profile.role == 'admin'">Редактировать</button>
              <button @click="deleteProduct(product.id)" class="btn btn-danger mt-2"  v-if="profile.role == 'admin'">Удалить</button> -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <footer>
    <nav aria-label="Page navigation">
      <ul class="pagination justify-content-center mt-4">
        <li class="page-item">
          <button class="page-link border-dark border-3 rounded-3 py-2 px-4" @click="changePage(1)"
            :disabled="page === 1">
            <i class="fas fa-fast-backward"></i> В начало
          </button>
        </li>

        <li class="page-item">
          <button class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4"
            @click="changePage(page - 1)" :disabled="page === 1">
            <i class="fas fa-chevron-left"></i> Назад
          </button>
        </li>

        <li class="page-item" v-for="pageNumber in visiblePages" :key="pageNumber"
          :class="{ active: page === pageNumber }">
          <button class="page-link border-dark border-3 rounded-3 py-2 px-4" @click="goToPage(pageNumber)"
            :disabled="page === pageNumber">
            {{ pageNumber }}
          </button>
        </li>

        <li class="page-item">
          <button class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4"
            @click="changePage(page + 1)" :disabled="page === totalPages">
            <i class="fas fa-chevron-right"></i> Вперед
          </button>
        </li>

        <li class="page-item">
          <button class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4"
            @click="changePage(totalPages)" :disabled="page === totalPages">
            В конец <i class="fas fa-fast-forward"></i>
          </button>
        </li>
      </ul>
    </nav>
  </footer>
</template>

<script>
import toastr from "toastr";
import "../../node_modules/toastr/build/toastr.min.css";
import axios from 'axios';
import { offset } from "@popperjs/core";

export default {
  data() {
    return {
      offset: 0,
      products: [], // Список товаров
      totalPages: null,  // Общее количество страниц
      page: 1, // Текущая страница
      limit: 9, // Количество товаров на страницу
      loading: false, // Флаг загрузки

      product: {
        name: '',
        description: '',
        price: 0,
        stock: 0,
        category: '',
        picture_url: ''
      },

      profile: null,
      errorMessage: '',
    };
  },
  computed: {
    visiblePages() {
      const maxPagesToShow = 5; // Максимум страниц для отображения
      let pages = [];

      // Стартовая и конечная страницы для отображения
      let startPage = Math.max(1, this.page - 2);
      let endPage = Math.min(this.totalPages, this.page + 2);

      // Если в начале не хватает страниц, сдвигаем правую сторону
      if (this.page <= 2) {
        endPage = Math.min(this.totalPages, maxPagesToShow);
      }

      // Если в конце не хватает страниц, сдвигаем левую сторону
      if (this.page >= this.totalPages - 2) {
        startPage = Math.max(1, this.totalPages - maxPagesToShow + 1);
      }

      // Заполняем массив номеров страниц
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i);
      }

      return pages;
    },
  },
  methods: {
    goToPage(pageNumber) {
      if (pageNumber >= 0 && pageNumber <= this.totalPages) {
        this.page = pageNumber
        this.fetchProducts();
      }
    },

    async fetchProfile() {
      try {
        const response = await axios.get("http://localhost:8080/profile", {
          withCredentials: true,
        });
        // console.log()
        this.profile = response.data;
      } catch (error) {
        this.errorMessage =
          error.response?.data?.message || "Ошибка загрузки профиля";
      }
    },

    async addToCart(productId) {
      try {
        var z = localStorage.getItem("user");

        if (z) {
          try {
            var parsedUser = JSON.parse(z);
            z = parsedUser.id || null; // Извлекаем id, если он есть
          } catch (e) {
            console.error("Ошибка при парсинге JSON:", e);
          }
        }
        console.log(z, productId);

        await axios.post("http://localhost:8080/add_to_cart", {
          user_id: z,
          product_id: productId,
          quantity: 1,
        });
        console.log(productId);
        toastr.success("Товар добавлен в корзину!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      } catch (error) {
        console.error("Ошибка при добавлении в корзину:", error);
        toastr.error("Не удалось добавить товар в корзину.", "Ошибка", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      }
    },
    async fetchProducts() {
      this.loading = true;
      this.limit = 9
      try {
        const response = await axios.get(`http://localhost:8080/products?page=${this.page}&limit=${this.limit}`
          , {
            params: { page: this.page, limit: this.limit },
          });

        this.products = response.data.products;
        this.totalPages = response.data.total_pages;
        this.page = response.data.page;
        console.log(response.data);

        // Обновляем URL с текущими параметрами
        this.updateURL();
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
      } finally {
        this.loading = false;
      }
      console.log(this.products)
    },
    updateURL() {
      // Синхронизация адресной строки
      this.$router.push({
        path: '/products',
        query: { page: this.page, limit: this.limit },
      });
    },
    changePage(newPage) {
      if (newPage > 0) {
        this.page = newPage;
        this.fetchProducts();
      }
    },
  },
  watch: {
    // Следим за изменением строки запроса
    $route(to) {
      const page = parseInt(to.query.page) || 1;
      const limit = parseInt(to.query.limit) || 2;

      if (page !== this.page || limit !== this.limit) {
        this.page = page;
        this.limit = limit;
        this.fetchProducts();
      }
    },
  },
  created() {
    const page = parseInt(this.$route.query.page) || 1;
    const limit = parseInt(this.$route.query.limit) || 2;
    this.page = page;
    this.limit = limit;

    this.fetchProducts();
    this.fetchProfile();
  },
};
</script>

<style scoped>
@import "../../node_modules/toastr/build/toastr.min.css";

.page-item button:disabled {
  background-color: #dcdcdc;
  color: #888888;
  cursor: not-allowed;
  border-color: #dcdcdc;
}

.page-item button:disabled i {
  color: #888888;
}


.toast {
  background-color: #28a745 !important;
  color: #fff !important;
}

.upload-box {
  padding: 20px;
  border: 2px dashed #007bff;
  text-align: center;
  border-radius: 8px;
  transition: background-color 0.3s ease;
}

.upload-box.drag-active {
  background-color: #f0f8ff;
}

.upload-link {
  color: #007bff;
  cursor: pointer;
  text-decoration: underline;
}

.upload-link:hover {
  text-decoration: none;
}

.upload-area {
  border: 2px dashed #ccc;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  position: relative;
}

.upload-area img {
  max-width: 100%;
  margin-top: 10px;
}

.error-message {
  color: red;
  margin-top: 10px;
}

.dropzone {
  width: 100%;
  height: 200px;
  border: 2px dashed #aaa;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.dropzone:hover {
  background-color: #f0f0f0;
}

button {
  margin-top: 20px;
}

.error-message,
.success-message {
  margin-top: 20px;
  padding: 10px;
  background-color: #f8d7da;
  border: 1px solid #f5c6cb;
  color: #721c24;
}

.success-message {
  background-color: #d4edda;
  border: 1px solid #c3e6cb;
  color: #155724;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 80%;
  max-width: 500px;
}

.form input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.card-img-top {
  object-fit: cover;
  height: 150px;
  max-height: 150px;
  width: 100%;
}

.card-body {
  display: flex;
  flex-direction: column;
}

.pagination-container {
  margin-top: 20px;
}

.card14 {
  height: 600px;
  max-height: 600px;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.card-text {
  max-height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #666;
  font-size: 0.95rem;
}

.card-text:last-child {
  font-weight: bold;
  color: #ff4b2b;
}

.row {
  gap: 20px;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.col-12.col-sm-6.col-md-4.col-lg-3 {
  display: flex;
  justify-content: center;
}

.pagination button {
  margin: 0 10px;
}

.btn {
  margin-top: auto;
}

.card-body {
  display: flex;
  flex-direction: column;
}

form {
  margin-bottom: 20px;
}

form div {
  margin-bottom: 10px;
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.product-card {
  margin-bottom: 20px;
}

.pagination button {
  margin: 0 10px;
}

.card {
  width: 100%;
  height: 350px;
  border: none;
  border-radius: 15px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  background: #ffffff;
  opacity: 0;
  transform: translateY(20px);
  animation: fadeIn 0.5s ease forwards;
}

.card:hover {
  transform: translateY(-10px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

.card-img-top {
  object-fit: cover;
  height: 150px;
  max-height: 150px;
  width: 100%;
  transition: transform 0.3s ease;
}

.card:hover .card-img-top {
  transform: scale(1.05);
}

.card:nth-child(odd) {
  animation-delay: 0.1s;
}

.card:nth-child(even) {
  animation-delay: 0.2s;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.card-text {
  max-height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #666;
  font-size: 0.95rem;
  margin-top: 10px;
}

.card-text:last-child {
  font-weight: bold;
  color: #ff4b2b;
}

.row {
  gap: 20px;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.col-12.col-sm-6.col-md-4.col-lg-3 {
  display: flex;
  justify-content: center;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.pagination-container {
  margin-bottom: 30px;
  z-index: 999;
  padding: 0px;
}

.card-img-top {
  width: 100%;
  height: 200px;
  object-fit: contain;
}
</style>