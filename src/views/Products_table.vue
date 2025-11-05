<template>

  <link href="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.css" rel="stylesheet">
  <div>
    <h1 style="display: flex; justify-content: center; margin: 30px;">Склад</h1>
    <button @click="openModal" style="margin-right: 50px; width: 100%;" class="btn btn-success mb-3">Добавить товар на
      склад</button>

    <button @click="printCart" style="width: 100%;" class="btn btn-primary mb-3">
      Скачать PDF с таблицей (Скачивает только те, что видны на странице) (с сервера)
    </button>

    <button @click="printTable" style="width: 100%;" class="btn btn-primary mb-3">
      Скачать PDF с таблицей (Скачивает только те, что видны на странице) (с клиента)
    </button>

    <div class="mt-3">
      <label for="productLimit">Введите количество продуктов, которые нужно отобразить (от 1 до -):</label>
      <input style="max-width: 100%;" type="number" id="productLimit" v-model="limit" class="form-control" min="1"
        step="1" />
    </div>

    <button @click="updateLimit(limit)" style="width: 100%;" class="btn btn-primary mt-3">Показать {{ limit }}
      товаров</button>
  </div>

  <!-- Модальное окно для добавления нового товара -->
  <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
    <div class="modal-content">
      <h5 class="modal-title">Добавить новый товар</h5>
      <form @submit.prevent="addProduct">
        <div class="form-group">
          <label for="productName">Название товара</label>
          <input type="text" id="productName" v-model="newProduct.Name" class="form-control" minlength="2"
            maxlength="99" required placeholder="Введите название товара" />
        </div>

        <div class="form-group">
          <label for="productDescription">Описание</label>
          <textarea id="productDescription" v-model="newProduct.Description" class="form-control" rows="3" minlength="2"
            maxlength="99" required placeholder="Введите описание товара"></textarea>
        </div>

        <div class="form-group">
          <label for="productPrice">Цена</label>
          <input type="number" id="productPrice" v-model="newProduct.Price" class="form-control" minlength="2"
            maxlength="99" required placeholder="Введите цену" min="0" step="1" max="999999999" />
        </div>

        <div class="form-group">
          <label for="productStock">Остаток</label>
          <input type="number" id="productStock" v-model="newProduct.Stock" class="form-control" required
            placeholder="Введите количество на складе" min="0" step="1" max="999999999" />
        </div>

        <div class="form-group">
          <label for="productCategory">Категория</label>
          <input type="text" id="productCategory" v-model="newProduct.Category" class="form-control" minlength="2"
            maxlength="99" required placeholder="Введите категорию товара" />
        </div>

        <div class="form-group">
          <label for="productPicture">Изображение</label>
          <input type="file" id="productPicture" @change="handleFileChange" class="form-control" accept="image/*" />
          <!-- Предпросмотр изображения -->
          <div v-if="newProduct.Picture_url" class="mt-3 text-center">
            <img :src="newProduct.Picture_url" alt="Предпросмотр изображения" class="preview-image" />
          </div>
        </div>

        <div class="form-group text-center">
          <button type="submit" class="btn btn-primary">Добавить товар</button>
          <button type="button" class="btn btn-secondary ml-2" @click="closeModal">
            Отмена
          </button>
        </div>
      </form>
    </div>
  </div>

  <!-- Модальное окно для редактирования товара -->
  <div v-if="isEditModalOpen" class="modal-overlay" @click.self="closeEditModal">
    <div class="modal-content">
      <h5 class="modal-title">Редактировать товар</h5>
      <form @submit.prevent="editProduct">
        <div class="form-group">
          <label for="editProductName">Название товара</label>
          <input type="text" id="editProductName" v-model="currentProduct.name" class="form-control" minlength="2"
            maxlength="99" required />
        </div>
        <div class="form-group">
          <label for="editProductDescription">Описание</label>
          <textarea id="editProductDescription" v-model="currentProduct.description" class="form-control" rows="3"
            minlength="2" maxlength="254" required></textarea>
        </div>
        <div class="form-group">
          <label for="editProductPrice">Цена</label>
          <input type="number" id="editProductPrice" v-model="currentProduct.price" class="form-control" min="0"
            max="999999999" required />
        </div>
        <div class="form-group">
          <label for="editProductStock">Остаток</label>
          <input type="number" id="editProductStock" v-model="currentProduct.stock" class="form-control" min="0"
            max="999999999" required />
        </div>
        <div class="form-group">
          <label for="editProductCategory">Категория</label>
          <input type="text" id="editProductCategory" v-model="currentProduct.category" class="form-control"
            minlength="2" maxlength="49" required />
        </div>
        <div class="form-group">
          <label for="editProductPicture">Изображение</label>
          <input type="file" id="editProductPicture" @change="handleEditFileChange" class="form-control" />
          <label for="editProductPicture">{{ currentProduct.Picture_url }}</label>
          <img v-if="currentProduct.Picture_url" :src="currentProduct.Picture_url" alt="Product Image"
            class="preview-image mt-2" />
        </div>
        <div class="form-group text-center">
          <button type="submit" class="btn btn-primary">Сохранить изменения</button>
          <button type="button" class="btn btn-secondary ml-2" @click="closeEditModal">Отмена</button>
        </div>
      </form>
    </div>
  </div>
  <div>

    <div v-if="loading">Загрузка...</div>
    <div v-else style="margin-top: 30px;">
      
      <div class="table-responsive">
        <table class="table table-striped table-bordered table-hover">
          <thead class="thead-dark">
            <tr>
              <th>№</th>
              <th>Название</th>
              <th>Описание</th>
              <th>Цена</th>
              <th>Изображение</th>
              <th>Остаток</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(product, index) in products" :key="product.id">
              <td>{{ (page - 1) * limit + index + 1 }}</td>
              <td>
                <router-link :to="'/product/' + product.id" class="text-decoration-none text-dark">
                  <strong>{{ product.name }}</strong>
                </router-link>
              </td>
              <td>{{ product.description }}</td>
              <td>{{ product.price }} ₽</td>
              <td>
                <img :src="product.Picture_url" alt="picture" class="product-image" />
              </td>
              <td>
                {{ product.stock }}
              </td>
              <td>
                <button @click="openEditModal(product)" style="width: 130px;"
                  class="btn btn-warning btn-sm mt-2">Редактировать</button>
              </td>
              <td>
                <button @click="deleteProduct(product.id)" style="width: 80px;"
                  class="btn btn-danger btn-sm mt-2">Удалить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <footer>
    <nav aria-label="Page navigation">
      <ul class="pagination justify-content-center mt-4">
        <!-- Кнопка "В начало" -->
        <li class="page-item">
          <button style="margin-right: 10px;" class="page-link border-dark border-3 rounded-3 py-2 px-4 button1" @click="changePage(1)"
            :disabled="page === 1">
            <i class="fas fa-fast-backward"></i> В начало
          </button>
        </li>

        <!-- Кнопка "Назад" -->
        <li class="page-item">
          <button style="margin-right: 10px;" class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4 button1"
            @click="changePage(page - 1)" :disabled="page === 1">
            <i class="fas fa-chevron-left"></i> Назад
          </button>
        </li>

        <!-- Кнопки с номерами страниц -->
        <li class="page-item" v-for="pageNumber in visiblePages" :key="pageNumber"
          :class="{ active: page === pageNumber }">
          <button class="page-link border-dark border-3 rounded-3 py-2 px-4 button1" @click="goToPage(pageNumber)"
            :disabled="page === pageNumber">
            {{ pageNumber }}
          </button>
        </li>

        <!-- Кнопка "Вперед" -->
        <li class="page-item">
          <button style="margin-left: 10px; margin-right: 10px;" class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4 button1"
            @click="changePage(page + 1)" :disabled="page === totalPages">
            <i class="fas fa-chevron-right"></i> Вперед
          </button>
        </li>

        <!-- Кнопка "В конец" -->
        <li class="page-item">
          <button style="margin-right: 10px;" class="page-link border-dark border-3 border-start-0 rounded-3 py-2 px-4 button1"
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
import axios from "axios";

export default {
  data() {
    return {

      totalPages: null,
      products: [], // Список товаров
      page: 1, // Текущая страница
      limit: 10, // Количество товаров на страницу
      loading: false, // Флаг загрузки

      isEditModalOpen: false,
      currentProduct: {
        id: null,
        name: "",
        description: "",
        price: 0,
        stock: 0,
        category: "",
        Picture_url: "",
        file: null,
      },

      isModalOpen: false, // Статус модального окна
      newProduct: { // Новый товар
        Name: '',
        Description: '',
        Price: '',
        Stock: '',
        Category: '',
        Picture_url: null,
      },
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
    }
  },
  methods: {
    printTable() {
    // Генерация HTML для таблицы с добавленными стилями
    const tableHtml = `
      <html>
      <head>
        <title>Печать таблицы</title>
        <style>
          body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
          table { width: 100%; border-collapse: collapse; margin-top: 20px; }
          th, td { padding: 12px; text-align: center; border: 1px solid #ddd; }
          th { background-color: #3498db; color: white; }
          td { background-color: #f9f9f9; }
          .product-image { width: 50px; height: 50px; object-fit: cover; }
          .thead-dark { background-color: #343a40; color: white; }
          .btn { padding: 5px 10px; font-size: 12px; cursor: pointer; }
          .table-responsive { margin: 20px; }
        </style>
      </head>
      <body>
        <div class="table-responsive">
          <table class="table table-striped table-bordered table-hover">
            <thead class="thead-dark">
              <tr>
                <th>№</th>
                <th>Название</th>
                <th>Описание</th>
                <th>Цена</th>
                <th>Изображение</th>
                <th>Остаток</th>
              </tr>
            </thead>
            <tbody>
              ${this.products.map((product, index) => `
                <tr>
                  <td>${(this.page - 1) * this.limit + index + 1}</td>
                  <td>
                    <a href="/product/${product.id}" class="text-decoration-none text-dark">
                      <strong>${product.name}</strong>
                    </a>
                  </td>
                  <td>${product.description}</td>
                  <td>${product.price} ₽</td>
                  <td>
                    <img src="${product.Picture_url}" alt="picture" class="product-image" />
                  </td>
                  <td>${product.stock}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      </body>
      </html>
    `;

    // Открываем новое окно для печати
    const printWindow = window.open('', '', 'width=800,height=600');

    // Заполняем окно HTML-контентом
    printWindow.document.write(tableHtml);
    printWindow.document.close();

    // Ждем загрузки содержимого и запускаем печать
    printWindow.onload = function () {
      printWindow.print();
  }
},
    // Метод для валидации поля
    validateField(fieldName, value, minLength, maxLength) {
      if (value.length < minLength) {
        this.errors[fieldName] = `Minimum length is ${minLength} characters.`;
      } else if (value.length > maxLength) {
        this.errors[fieldName] = `Maximum length is ${maxLength} characters.`;
      } else {
        this.errors[fieldName] = '';
      }
    },

    goToPage(pageNumber) {
      if (pageNumber >= 0 && pageNumber <= this.totalPages) {
        this.page = pageNumber
        this.fetchProducts();
      }
    },

    openEditModal(product) {
      this.isEditModalOpen = true;
      this.currentProduct = { ...product }; // Копируем данные
      console.log(product)
      console.log(this.currentProduct)
    },

    closeEditModal() {
      this.isEditModalOpen = false;
    },

    updateLimit(inputLimit) {
      console.log(`Текущий лимит: ${inputLimit}`);
      this.$router.push({ query: { limit: inputLimit } });
      this.fetchProducts(inputLimit);
    },

    handleEditFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        this.currentProduct.Picture_url = URL.createObjectURL(file);
        this.currentProduct.file = file;
      } else {
        this.currentProduct.Picture_url = null;
        this.currentProduct.file = null;
      }
    },

    async editProduct() {
      try {
        const formData = new FormData();
        formData.append("Name", this.currentProduct.name);
        formData.append("Description", this.currentProduct.description);
        formData.append("Price", this.currentProduct.price);
        formData.append("Stock", this.currentProduct.stock);
        formData.append("Category", this.currentProduct.category);

        if (this.currentProduct.file) {
          formData.append("Picture_url", this.currentProduct.file);
        }

        const response = await fetch(`http://localhost:8080/product_edit/${this.currentProduct.id}`, {
          method: "POST",
          body: formData,
        });

        if (!response.ok) {
          throw new Error("Ошибка при обновлении товара");
        }

        const data = await response.json();
        console.log("Обновленный товар:", data);

        const index = this.products.findIndex((p) => p.id === this.currentProduct.id);
        if (index !== -1) {
          this.products.splice(index, 1, data.product);
        }

        this.closeEditModal();
        toastr.success("Товар изменен на складе!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
        window.location.reload();

      } catch (error) {
        console.error("Ошибка:", error.message);
        toastr.error("Не удалось отредактировать товар.", "Ошибка", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      }
    },

    async deleteProduct(productId) {
      try {
        await axios.delete(`http://localhost:8080/products_delete/${productId}`);
        alert('Товар удален!');
        this.fetchProducts();
        toastr.success("Товар удален со склада!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      } catch (error) {
        console.error('Ошибка при удалении товара:', error);

        toastr.error("Не удалось удалить товар.", "Ошибка", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      }
    },

    openModal() {
      this.isModalOpen = true;
    },

    closeModal() {
      this.isModalOpen = false;
      this.resetNewProduct();
    },

    resetNewProduct() {
      this.newProduct = {
        name: '',
        description: '',
        price: '',
        stock: '',

        category: '',
        pictureUrl: null,
      };
    },

    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        this.newProduct.Picture_url = URL.createObjectURL(file); // Для предпросмотра
        this.newProduct.file = file;
      } else {
        this.newProduct.Picture_url = null;
        this.newProduct.file = null;
      }
    },

    async addProduct() {
      try {
        const formData = new FormData();
        formData.append('Name', this.newProduct.Name);
        formData.append('Description', this.newProduct.Description);
        formData.append('Price', this.newProduct.Price);
        formData.append('Stock', this.newProduct.Stock);
        formData.append('Category', this.newProduct.Category);

        if (this.newProduct.file) {
          formData.append('Picture_url', this.newProduct.file); // Добавляем файл только если он выбран
        }

        const response = await fetch('http://localhost:8080/products_add', {
          method: 'POST',
          body: formData,
        });

        if (!response.ok) {
          throw new Error('Ошибка при добавлении товара');
        }

        const data = await response.json();
        console.log('Добавленный товар:', data);
        this.products.push(data.product); // Добавляем товар в локальный список
        this.closeModal();

        toastr.success("Товар добавлен на склад!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
        window.location.reload();

      } catch (error) {
        console.error('Ошибка:', error.message);

        toastr.error("Не удалось добавить товар.", "Ошибка", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      }
    },

    printCart() {
      const cartData = this.products.map(item => ({
        id: item.id,
        name: item.name,
        price: item.price,
        stock: item.stock,
        category: item.category
      }));
      console.log('ПЕЧАТАЮ' + cartData[1].price);

      fetch('http://localhost:8080/generate-pdf15', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ cart: cartData })
      })
        .then(response => response.blob())
        .then(blob => {
          const link = document.createElement('a');
          link.href = URL.createObjectURL(blob);
          link.download = 'cart.pdf';
          link.click();
        })
        .catch(error => console.error('Error:', error));

      // Показываем уведомление Toastr
      toastr.success("Файл скачан!", "Успех", {
        closeButton: true,
        progressBar: true,
        positionClass: "toast-bottom-right",
      });
    },

    async fetchProducts() {
      this.loading = true;
      try {
        const response = await axios.get(`http://localhost:8080/products_table?page=${this.page}&limit=${this.limit}`
          , {
            params: { page: this.page, limit: this.limit },
          });

        this.products = response.data.products;
        this.totalPages = response.data.total_pages; // Общее число страниц
        this.page = response.data.page; // Текущая страница
        console.log(response.data);

        // Обновляем URL с текущими параметрами
        this.updateURL();
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
      } finally {
        this.loading = false;
      }
    },
    updateURL() {
      // Синхронизация адресной строки
      this.$router.push({
        path: '/products_table',
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
  },
};
</script>

<style scoped>
.button1 {
  margin-right: 10px;
  margin-top: 20px;
}

.button1 {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

.button1:hover {
  background-color: #0056b3;
}

.page-item button:disabled {
  background-color: #dcdcdc;
  color: #888888;
  cursor: not-allowed;
  border-color: #dcdcdc;
}

.page-item button:disabled i {
  color: #888888;
}

/* Модальное окно */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1050;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  max-width: 500px;
  width: 100%;
}

.modal-title {
  font-size: 1.25rem;
  margin-bottom: 20px;
  text-align: center;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  font-weight: bold;
}

.preview-image {
  max-width: 100%;
  max-height: 200px;
  margin-top: 10px;
  object-fit: cover;
  border-radius: 5px;
}

.btn {
  width: 48%;
}

.btn-secondary {
  background-color: #6c757d;
}

.btn-primary {
  background-color: #007bff;
}

.pagination-container {

  margin-top: 20px;

  margin-bottom: 30px;
  z-index: 999;
  padding: 0px;
}

.pagination button {

  max-width: 130px;
  margin: 0px;
}

.product-image {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 5px;
}

.table th,
.table td {
  vertical-align: middle;
}

.table-striped tbody tr:nth-of-type(odd) {
  background-color: #f8f9fa;
}

.table-hover tbody tr:hover {
  background-color: #e9ecef;
}

th {
  text-align: center;
  font-weight: bold;
}

td {
  text-align: center;
}

.btn {
  font-size: 0.875rem;
  padding: 0.375rem 0.75rem;
  border-radius: 0.375rem;
}

.btn-warning {
  background-color: #ffc107;
  border-color: #ffc107;
}

.btn-danger {
  background-color: #dc3545;
  border-color: #dc3545;
}

.product-image {
  width: 50px;
  height: 50px;
  object-fit: cover;
}

.table {
  margin-top: 20px;
}

.product-image {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 5px;
}
</style>