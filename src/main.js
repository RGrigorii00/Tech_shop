// import './assets/main.css'

import 'aos/dist/aos.css';  // Импортируем стили AOS
import AOS from 'aos'; // Импортируем саму библиотеку

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap';
import * as THREE from 'three';



AOS.init(); // Инициализируем AOS

const app = createApp(App)

app.use(router)

app.mount('#app')
