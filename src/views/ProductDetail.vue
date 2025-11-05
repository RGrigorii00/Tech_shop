<template>

  <link href="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.css" rel="stylesheet">
  <div class="container mt-5">
    <div class="row">
      <!-- Верхний раздел -->
      <div class="col-12">
        <div class="row">
          <!-- 3D модель слева -->
          <div v-if="this.osType == 'Windows' || this.osType == 'MacOS'" class="col-12 col-md-6">
            <h1 style="display: flex; justify-content: center;">3D модель</h1>
            <div class="threejs-container" ref="modelContainer" style="width: 600px; height: 600px;"></div>
          </div>

          <!-- Картинки товара справа -->
          <div v-if="this.isReady == true" class="col-12 col-md-6">
            <div id="product-carousel" class="carousel slide" data-bs-ride="carousel">
              <h1 style="display: flex; justify-content: center;">Картинки товара</h1>
              <div class="carousel-inner">
                <div class="carousel-item active">
                  <img :src="this.product.picture_url" class="d-block w-100 img-fluid abc" alt="Товар 1" />
                </div>
                <div class="carousel-item">
                  <img :src="this.product.picture_url" class="d-block w-100 img-fluid abc" alt="Товар 2" />
                </div>
                <div class="carousel-item">
                  <img :src="this.product.picture_url" class="d-block w-100 img-fluid abc" alt="Товар 3" />
                </div>
              </div>

              <!-- Стрелочки для перелистывания -->
              <button class="carousel-control-prev" type="button" data-bs-target="#product-carousel"
                data-bs-slide="prev">
                <span class="carousel-control-prev-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Предыдущий</span>
              </button>
              <button class="carousel-control-next" type="button" data-bs-target="#product-carousel"
                data-bs-slide="next">
                <span class="carousel-control-next-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Следующий</span>
              </button>
            </div>

            <!-- Миниатюры картинок -->
            <div class="d-flex mt-3">
              <button type="button" data-bs-target="#product-carousel" data-bs-slide-to="0" class="thumb-item me-2">
                <img :src="this.product.picture_url" class="img-thumbnail w-100 h-100" alt="Товар 1" />
              </button>
              <button type="button" data-bs-target="#product-carousel" data-bs-slide-to="1" class="thumb-item me-2">
                <img :src="this.product.picture_url" class="img-thumbnail w-100 h-100" alt="Товар 2" />
              </button>
              <button type="button" data-bs-target="#product-carousel" data-bs-slide-to="2" class="thumb-item">
                <img :src="this.product.picture_url" class="img-thumbnail w-100 h-100" alt="Товар 3" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Нижний раздел -->
      <div class="col-12">
        <div class="row">
          <!-- Параметры модели слева -->
          <div class="col-12 col-md-6">
            <div v-if="this.osType == 'Windows' || this.osType == 'MacOS'" class="mt-5">

              <h3 style="display: flex; justify-content: center;">Параметры 3D модели</h3>
              <h6 style="display: flex; justify-content: center;">Система - {{ this.osType }}
                (На Windows, MaсOS
                работает. На Linux(Android) отключена 3D модель и ее параметры).
                Колесико - приблезить/отдалить. ЛКМ - вращать. ПКМ - перемещать.
                Модельки с текстурами есть только у wd, клавиатуры и телефона, в остальных - заглушка.</h6>
              <button @click="randomizeModel">Randomize Model</button>
              <div>
                <label>Light Type:</label>
                <select v-model="lightType" @change="updateLightType">
                  <option value="0">Ambient & Directional</option>
                  <option value="1">Point Light</option>
                  <option value="2">Colored Light</option>
                </select>
              </div>

              <div>
                <label>Light Intensity:</label>
                <input type="range" v-model="lightIntensity" min="0" max="2" step="0.1" @input="updateLightIntensity" />
                <span>{{ lightIntensity }}</span>
              </div>

              <div>
                <label>Light Color:</label>
                <input type="color" v-model="lightColor" @input="updateLightColor" />
              </div>

              <div>
                <label>Model Size (X):</label>
                <input type="range" v-model="modelSizeX" min="0.1" max="3" step="0.1" @input="updateModelSizeX" />
                <span>{{ modelSizeX }}</span>
              </div>

              <div>
                <label>Model Size (Y):</label>
                <input type="range" v-model="modelSizeY" min="0.1" max="3" step="0.1" @input="updateModelSizeY" />
                <span>{{ modelSizeY }}</span>
              </div>

              <div>
                <label>Model Size (Z):</label>
                <input type="range" v-model="modelSizeZ" min="0.1" max="3" step="0.1" @input="updateModelSizeZ" />
                <span>{{ modelSizeZ }}</span>
              </div>
            </div>
          </div>

          <!-- Информация о товаре справа -->
          <div class="col-12 col-md-6">
            <div class="product-details mt-4">
              <h2 class="product-title">{{ this.product.name }}</h2>
              <p class="product-price">Цена: {{ this.product.price }} ₽</p>
              <ul class="list-group">
                <li class="list-group-item"><strong>Описание: </strong> {{ this.product.description }}</li>
                <li class="list-group-item"><strong>Осталось:</strong> {{ this.product.stock }} товар</li>
                <li class="list-group-item"><strong>Категория:</strong> {{ this.product.category }}</li>
                <li class="list-group-item"><strong>Привезен:</strong> {{ formatDateTimeFromString(this.product.created_at) }}</li>
                <li class="list-group-item"><strong>Картинка:</strong> {{ this.product.picture_url }}</li>
              </ul>
              <button style="margin-right: 10px;" class="btn btn-primary mt-3 mg-5" @click="addToCart(this.product.id)">Добавить в корзину</button>

              <button class="btn btn-primary mt-3" @click="print1">Скачать PDF</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import toastr from "toastr";
import axios from 'axios';
import ThreeModel from '../../3dmodels/threejsModel.js'; // Импортируем класс ThreeModel

export default {
  data() {
    return {

      isReady: false,

      product: {},
      Phoro_url: null,
      profile: null,
      role: null,

      osType: this.detectOS(),
      threeModel: null, // Экземпляр ThreeModel
      // modelUrl: '../../3dmodels/3dmodels/gaming_keyboard.glb', // Путь к 3D-модели
      // textureUrl: '../../3dmodels/textures/wd/wd1tb-green_baseColor.png', // Путь к текстуре

      lightType: 0, // Тип света (0, 1, 2)
      lightIntensity: 1, // Интенсивность света
      lightColor: "#ffffff", // Цвет света

      modelSizeX: 1, // Размер модели по оси X
      modelSizeY: 1, // Размер модели по оси Y
      modelSizeZ: 1, // Размер модели по оси Z
    };
  },
  mounted() {},
  beforeUnmount() {
    window.removeEventListener("resize", this.handleResize);
  },
  created() {
    // Сначала загружаем товар, потом модель
    this.fetchProduct().then(() => {
      this.initThreeModel();
      this.fetchProfile();
    });
    window.addEventListener("resize", this.handleResize);
  },
  methods: {
    async fetchProfile() {
      try {
        const response = await axios.get("http://localhost:8080/profile15", {
          withCredentials: true,
        });

        console.log(response.data.user.photo_url);
        this.profile = response.data.user;
        console.log(this.profile)
      } catch (error) {
        this.errorMessage =
          error.response?.data?.message || "Ошибка загрузки профиля";
  }
},
    print1() {
      const printContent = `
<html>
  <head>
    <title>Карточка товара</title>
    <style>
      /* Настройка страницы для печати */
      @page {
        size: A4 landscape; /* Горизонтальная ориентация */
        margin: 20mm; /* Устанавливаем отступы */
      }

      body {
        font-family: 'Arial', sans-serif;
        margin: 0;
        padding: 0;
        background-color: #f9f9f9;
        color: #333;
        font-size: 14px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        width: 100%;
      }

      h1 {
        text-align: center;
        font-size: 28px;
        color: #4a90e2;
        margin-bottom: 20px;
        font-weight: bold;
      }

      .print-container {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        width: 100%;
        max-width: 900px;
        margin-top: 20px;
      }

      .image-section {
        width: 350px; /* Увеличена ширина изображения */
        height: auto;
        margin-right: 30px;
        text-align: center;
      }

      .info-section {
        width: calc(100% - 380px); /* Уменьшена ширина секции с текстом */
        padding-left: 30px;
        display: flex;
        flex-direction: column;
      }

      .info-section p {
        font-size: 16px;
        line-height: 1.6;
        color: #333;
        margin-bottom: 12px;
      }

      .product-name {
        font-size: 18px;
        font-weight: bold;
        color: #333;
        margin-bottom: 10px;
      }

      .product-description {
        font-style: italic;
        color: #7f8c8d;
        margin-bottom: 10px;
      }

      .product-price {
        font-size: 18px;
        color: #27ae60;
        margin-bottom: 10px;
      }

      .footer {
        text-align: center;
        font-size: 12px;
        color: #888;
        margin-top: 30px;
      }

      .print-image {
        width: 100%;
        height: auto;
        border-radius: 8px;
        object-fit: cover;
      }

      /* Печать */
      @media print {
        body {
          background-color: #ffffff;
          margin: 0;
          padding: 0;
          display: block;
        }

        .print-container {
          display: flex;
          justify-content: flex-start;
          align-items: flex-start;
          width: 100%;
          max-width: 100%;
          padding: 20px;
          box-shadow: none;
        }

        .image-section {
          width: 300px; /* Увеличена ширина изображения для печати */
        }

        .info-section {
          width: calc(100% - 340px);
        }

        .footer {
          font-size: 10px;
          margin-top: 40px;
        }

        h1 {
          font-size: 24px;
        }
      }
    </style>
  </head>
  <body>
    <div class="container">
      <h1>Карточка товара</h1>
      
      <div class="print-container">
        <!-- Секция с картинкой -->
        <div class="image-section">
          <img src="${this.product.picture_url}" alt="Картинка товара" class="print-image" />
        </div>

        <!-- Секция с информацией о товаре -->
        <div class="info-section">
          <p class="product-name"><strong>Название:</strong> ${this.product.name}</p>
          <p class="product-description"><strong>Описание:</strong> ${this.product.description}</p>
          <p class="product-price"><strong>Цена:</strong> ${this.product.price}</p>
          <p><strong>Осталось на складе:</strong> ${this.product.stock}</p>
          <p><strong>Категория:</strong> ${this.product.category}</p>
        </div>
      </div>

      <div class="footer">
        <p>© 2025 TechStore Рожков Григорий Олегович. Все права не защищены((((.</p>
      </div>
    </div>
  </body>
</html>
`;
      // Открываем новое окно для печати
      const printWindow = window.open('', '', 'width=800,height=600');

      // Заполняем окно созданным HTML
      printWindow.document.write(printContent);
      printWindow.document.close();
      printWindow.onload = function () {
        printWindow.print();
      };
    },

    initThreeModel() {
      if (this.product) {
        // Создаем экземпляр ThreeModel и передаем параметры
        this.threeModel = new ThreeModel(
          this.$refs.modelContainer,  // Контейнер для рендеринга
          this.product.model_url,      // URL модели
          this.product.texture_url,    // URL текстуры модели
          0                            // Мод по умолчанию для света (например, 0 для мягкого света)
        );
      }
    },

    async addToCart(productId) {
      try {
        var z = localStorage.getItem("user");
        if (z) {
          try {
            var parsedUser = JSON.parse(z);
            z = parsedUser.id || null;
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

formatDateTimeFromString(dateTimeString) {
  // Проверка на валидность входной строки
  if (!dateTimeString || typeof dateTimeString !== 'string') {
    throw new Error('Invalid input: Expected a string');
  }

  // Убираем символ "T" и разделяем строку на дату и время
  const [date, timeWithMilliseconds] = dateTimeString.split('T');

  if (!date || !timeWithMilliseconds) {
    throw new Error('Invalid dateTimeString format');
  }

  // Убираем возможный суффикс "Z" и миллисекунды
  const time = timeWithMilliseconds.split('.')[0]; // Оставляем только HH:MM:SS

  // Проверка валидности даты и времени
  if (!/^\d{4}-\d{2}-\d{2}$/.test(date)) {
    throw new Error('Invalid date format');
  }
  if (!/^\d{2}:\d{2}:\d{2}$/.test(time)) {
    throw new Error('Invalid time format');
  }

  // Форматирование: дата уже в формате YYYY-MM-DD, добавляем время
  return `${date} ${time}`;
},

    async fetchProduct() {
      try {
        const id = this.$route.params.id;
        const response = await axios.get(`http://localhost:8080/product/${id}`);
        const product = response.data;

        // Значения по умолчанию для model_url и texture_url
        const defaultModelUrl = '../../3dmodels/Chair.glb';
        const defaultTextureUrl = '../../3dmodels/normalsidyshkatexture_png.png';

        // Проверяем и заменяем пустые поля на значения по умолчанию
        product.model_url = product.model_url || defaultModelUrl;
        product.texture_url = product.texture_url || defaultTextureUrl;

        this.product = product;
        this.isReady = true;
        console.log(product)
      } catch (error) {
        console.error('Ошибка при получении товара:', error);
      }
    },

    getUserRole() {
      const i = localStorage.getItem('user');
      const parsedUser = JSON.parse(i);
      this.Phoro_url = parsedUser.phoro_url
      console.log("ЮЗЫР", parsedUser.phoro_url)
    },

    detectOS() {
      const userAgent = navigator.userAgent;
      if (userAgent.indexOf("Win") !== -1) {
        return "Windows";
      } else if (userAgent.indexOf("Mac") !== -1) {
        return "MacOS";
      } else if (userAgent.indexOf("Linux") !== -1) {
        return "Linux";
      } else if (userAgent.indexOf("Android") !== -1) {
        return "Android";
      } else if (userAgent.indexOf("iPhone") !== -1) {
        return "iOS";
      } else {
        return "Неизвестная ОС";
      }
    },

    // Метод для вызова randomizeAll
    randomizeModel() {
      if (this.threeModel) {
        this.threeModel.randomizeAll();
      }
    },

    // Обновление размера модели по оси X
    updateModelSizeX() {
      if (this.threeModel && this.threeModel.model) {
        this.threeModel.model.scale.x = parseFloat(this.modelSizeX);
      }
    },

    // Обновление размера модели по оси Y
    updateModelSizeY() {
      if (this.threeModel && this.threeModel.model) {
        this.threeModel.model.scale.y = parseFloat(this.modelSizeY);
      }
    },

    // Обновление размера модели по оси Z
    updateModelSizeZ() {
      if (this.threeModel && this.threeModel.model) {
        this.threeModel.model.scale.z = parseFloat(this.modelSizeZ);
      }
    },

    // Обновление типа света
    updateLightType() {
      if (this.threeModel) {
        this.threeModel.updateLightType(parseInt(this.lightType));
      }
    },

    // Обновление интенсивности света
    updateLightIntensity() {
      if (this.threeModel) {
        this.threeModel.updateLightIntensity(parseFloat(this.lightIntensity));
      }
    },

    // Обновление цвета света
    updateLightColor() {
      if (this.threeModel) {
        this.threeModel.updateLightColor(this.lightColor);
      }
    },
    handleResize() {
      if (this.threeModel) {
        this.threeModel.resize();
      }
    },
  },

  watch: {
    product(newProduct) {
      if (newProduct && this.threeModel) {
        this.threeModel.Model_url = newProduct.model_url;
        this.threeModel.Texture_url = newProduct.texture_url;
        this.threeModel.loadModel(); // Перезагружаем модель при изменении URL
      }
    }
  }
};
</script>

<style scoped>
.abc {
  /* object-fit: cover; */
  max-width: 600px;
  max-height: 600px;
  width: 600px;
  height: 600px;
}

.product-title {
  font-size: 2rem;
  font-weight: bold;
}

.product-price {
  font-size: 1.5rem;
  color: #007bff;
  margin-bottom: 1rem;
}

.list-group-item {
  font-size: 1.1rem;
  padding: 10px;
  border: 1px solid #ddd;
  transition: background-color 0.3s ease;
}

.list-group-item:hover {
  background-color: #f8f9fa;
}

.btn {
  transition: transform 0.3s;
}

.btn:hover {
  transform: scale(1.1);
}

.thumb-item img {
  max-width: 80px;
  max-height: 80px;
  object-fit: cover;
}
</style>