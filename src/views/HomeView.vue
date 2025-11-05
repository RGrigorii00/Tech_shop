<template>
  <div class="container py-5">
    <header class="text-center mb-5" data-aos="fade-down" data-aos-duration="1000">
      <h1 class="display-4">Добро пожаловать в TechStore, {{ userName }}!</h1>
      <p class="lead">Магазин техники, где технологии встречаются с качеством!</p>
      <button @click="viewProducts" class="btn btn-primary btn-lg">Посмотреть товары</button>
    </header>

    <section class="featured-products mt-5" data-aos="fade-up" data-aos-duration="1000">
      <h2 class="text-center mb-4">Рекомендуемые товары</h2>
      <div class="row">
        <div class="col-md-4 mb-4" data-aos="flip-left" data-aos-duration="1000">
          <div class="card">
            <img src="../assets/products_photos/SAMSUNG_phone.png" class="card-img-top" alt="Product 1">
            <div class="card-body">
              <h5 class="card-title">Смартфон XPro</h5>
              <p class="card-text">Новейшая модель смартфона с отличной камерой и большим экраном.</p>
              <p><strong>Цена: 40 000 ₽</strong></p>
              <button class="btn btn-success">Купить</button>
            </div>
          </div>
        </div>

        <div class="col-md-4 mb-4" data-aos="flip-left" data-aos-duration="1000">
          <div class="card">
            <img src="../assets/products_photos/ASUSROG14.png" class="card-img-top" alt="Product 2">
            <div class="card-body">
              <h5 class="card-title">Ноутбук ProBook</h5>
              <p class="card-text">Идеальный ноутбук для работы и учебы, с отличными характеристиками.</p>
              <p><strong>Цена: 60 000 ₽</strong></p>
              <button class="btn btn-success">Купить</button>
            </div>
          </div>
        </div>

        <div class="col-md-4 mb-4" data-aos="flip-left" data-aos-duration="1000">
          <div class="card">
            <img src="../assets/products_photos/Keyboard_asus_rog.png" class="card-img-top" alt="Product 3">
            <div class="card-body">
              <h5 class="card-title">Клавиатура</h5>
              <p class="card-text">Без проводов, с отличным качеством и длительным временем работы.</p>
              <p><strong>Цена: 5 000 ₽</strong></p>
              <button class="btn btn-success">Купить</button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="discounts mt-5" data-aos="fade-right" data-aos-duration="1000">
      <h2 class="text-center mb-4">Акции и скидки</h2>
      <div class="row">
        <div class="col-md-6 mb-4" data-aos="zoom-in" data-aos-duration="1000">
          <div class="card">
            <img src="../assets/products_photos/discount1.png" class="card-img-top" alt="Discount 1">
            <div class="card-body">
              <h5 class="card-title">Скидка 20% на смартфоны</h5>
              <p>Только до конца месяца! Получите 20% скидки на все модели смартфонов.</p>
              <button class="btn btn-danger">Узнать больше</button>
            </div>
          </div>
        </div>

        <div class="col-md-6 mb-4" data-aos="zoom-in" data-aos-duration="1000">
          <div class="card">
            <img src="../../src/assets/products_photos/discount1.png" class="card-img-top" alt="Discount 2">
            <div class="card-body">
              <h5 class="card-title">Скидка 10% на ноутбуки</h5>
              <p>Получите 10% скидку на все ноутбуки с кодом PROMO10.</p>
              <button class="btn btn-danger">Узнать больше</button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="about-us mt-5" data-aos="fade-left" data-aos-duration="1000">
      <h2 class="text-center mb-4">О магазине</h2>
      <p class="text-center">
        Мы — ваш надежный магазин техники. Мы предоставляем лучшие товары по выгодным ценам и обеспечиваем высокий
        уровень обслуживания. Наши клиенты могут всегда рассчитывать на качественную консультацию и поддержку.
      </p>
    </section>

    <section class="contact mt-5" data-aos="fade-up" data-aos-duration="1000">
      <h2 class="text-center mb-4">Контактная информация</h2>
      <ul class="list-unstyled text-center">
        <li><strong>Email:</strong> support@techstore.com</li>
        <li><strong>Телефон:</strong> +7 (123) 456-7890</li>
      </ul>
    </section>

    <footer class="text-center mt-5 py-3 bg-light" data-aos="" data-aos-duration="1000">
      <p>&copy; 2025 TechStore Рожков Григорий Олегович. Все права не защищены((((.</p>
    </footer>
  </div>
</template>

<script>
export default {
  name: "HomePage",
  data() {
    return {
      userName: null,
      isLight: true,
      showContactForm: false,
      form: {
        name: '',
        email: '',
        message: ''
      },
    };
  },

  created() {
    const user = JSON.parse(localStorage.getItem('user'));
    if (user && user.name) {
      this.userName = user.name;
    } else {
      this.userName = 'Гость';
    }
  },
  methods: {
    toggleBackground() {
      this.isLight = !this.isLight;
    },
    viewProducts() {
      this.$router.push({ name: 'products' });
    }
  }
};
</script>

<style scoped>
.card-img-top {
  object-fit: cover;
  height: 300px;
  max-height: 700px;
  max-width: 500px;
  width: 100%;
}

.card {
  border-radius: 10px;
}

.card-header {
  font-weight: bold;
}

footer {
  font-size: 14px;
}

.bg-light {
  background-color: #f8f9fa !important;
}

.bg-dark {
  background-color: #343a40 !important;
}

.card-img-top {
  border-radius: 10px;
}

.card-body {
  background-color: #f5f5f5;
}

.card-title {
  font-weight: bold;
}

.card-text {
  font-size: 16px;
}

h2 {
  font-size: 2rem;
}
</style>