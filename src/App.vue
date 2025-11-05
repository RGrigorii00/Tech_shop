<script>

import toastr from 'toastr';
import 'toastr/build/toastr.min.css';
import { RouterLink, RouterView } from 'vue-router'
import VuePaginate from 'vuejs-paginate'
import * as THREE from 'three';
</script>
<template>
    <header>
        <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <div class="container-fluid">
                <a class="navbar-brand" href="/"><img src="\src\assets\Logo.png" alt="" width="50" height="50"></a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
                    aria-controls="navbarNav" aria-expanded="false" aria-label="Togglenavigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav">
                        <li class="nav-item">
                            <router-link class="nav-link" to="/"><i class="bi-house"></i> Главная</router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link" to="/products"><i class="bi bi-bandaid"></i>
                                Товары</router-link>
                        </li>
                        <li class="nav-item" v-if="isLoggedIn && profile.role == 'admin'">
                            <router-link class="nav-link" to="/products_table">
                                Склад
                            </router-link>
                        </li>
                        <li class="nav-item">
                            <router-link class="nav-link" to="/about"><i class="bi bi-book"></i> О нас</router-link>
                        </li>
                        <li class="nav-item" v-if="isLoggedIn == true">
                            <router-link class="nav-link" to="/buy"><i class="bi bi-book"></i> Мои покупки</router-link>
                        </li>
                        <li class="nav-item" v-if="isLoggedIn && profile.role == 'admin'">
                            <router-link class="nav-link" to="/sells"><i class="bi bi-book"></i> Продажи</router-link>
                        </li>
                    </ul>
                </div>
                <div class="collapse navbar-collapse myclass" id="navbarNav">

                    <div style="margin-right: 15px;">
                        <router-link class="nav-link" to="/get_cart"><i class="bi bi-book"></i> Корзина</router-link>
                    </div>
                    <ul class="navbar-nav">
                        <!-- Кнопка Войти (показывается, если пользователь не залогинен) -->
                        <li class="nav-item" v-if="!isLoggedIn">
                            <router-link class="nav-link" to="/login"><i class="bi bi-book"></i>Войти</router-link>
                        </li>

                        <!-- Кнопка Регистрация (показывается, если пользователь не залогинен) -->
                        <li class="nav-item" v-if="!isLoggedIn">
                            <router-link class="nav-link" to="/register"><i
                                    class="bi bi-book"></i>Регистрация</router-link>
                        </li>

                        <!-- Кнопка Профиль (показывается, если пользователь залогинен) -->
                        <li class="nav-item" v-if="isLoggedIn">
                            <router-link class="nav-link" to="/profile">
                                <span v-if="profile" style="margin-right: 10px;">{{ profile.username }}</span>
                                <img v-if="profile" :src="profile.photo_url" class="user-avatar">
                        <li class="list-group-item" v-if="profile">
                            <strong>Баланс: </strong> {{ profile.balance }}₽
                        </li>
                        </router-link>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    </header>
    <main class="container">
        <RouterView />
    </main>
    <footer>
    </footer>
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import VuePaginate from 'vuejs-paginate'
import * as THREE from 'three';
import { ref, onMounted } from 'vue';
import axios from 'axios';

const isLoggedIn = ref(false);
const user = ref({
    avatar: '',
    name: ''
});
const profile = ref(null);
const errorMessage = ref('');
const cancelTokenSource = ref(null);
const isFetching = ref(false);

const checkLoginStatus = () => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
        isLoggedIn.value = true;
        user.value = JSON.parse(storedUser);
    }
};

const fetchProfile = async () => {
    if (cancelTokenSource.value) {
        cancelTokenSource.value.cancel('Новый запрос отменяет старый');
    }

    cancelTokenSource.value = axios.CancelToken.source();
    isFetching.value = true;

    try {
        const response = await axios.get('http://localhost:8080/profile', {
            withCredentials: true,
            cancelToken: cancelTokenSource.value.token,
        });

        console.log(response.data.user);
        profile.value = response.data.user;
    } catch (error) {
        if (axios.isCancel(error)) {
            console.log('Запрос отменен:', error.message);
        } else {
            errorMessage.value = error.response?.data?.message || 'Ошибка загрузки профиля';
        }
    } finally {
        isFetching.value = false;
    }
};

onMounted(() => {
    checkLoginStatus();
    fetchProfile();
});
</script>

<style scoped>
.user-avatar {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    margin-right: 10px;
}
</style>
<style scoped>
header {
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    color: #333333;
    box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    position: sticky;
    top: 0;
    z-index: 1000;
    transition: all 0.3s ease-in-out;
}

.navbar {
    background: transparent;
}

.navbar-brand img {
    filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.2));
    transition: transform 0.3s ease, filter 0.3s ease;
}

.navbar-brand img:hover {
    transform: scale(1.15);
    filter: drop-shadow(0 6px 8px rgba(254, 184, 79, 0.6));
}

.nav-link {
    color: rgba(51, 51, 51, 0.8);
    font-size: 1.2rem;
    font-weight: 500;
    letter-spacing: 0.5px;
    position: relative;
    transition: color 0.3s ease, transform 0.3s ease;
}

.nav-link::after {
    content: '';
    position: absolute;
    width: 0;
    height: 3px;
    background: linear-gradient(90deg, #fec74f, #feba00);
    bottom: -6px;
    left: 0;
    transition: width 0.4s ease;
    border-radius: 3px;
}

.nav-link:hover {
    color: #fea500;
    transform: translateY(-3px);
}

.nav-link:hover::after {
    width: 100%;
}

.collapse.show {
    animation: fadeIn 0.5s ease forwards;
}

.myclass {
    max-width: 300px;
    margin-right: 30px;
}
</style>
