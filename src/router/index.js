import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Products from '../views/Products.vue'
import ProductDetail from '@/views/ProductDetail.vue'

import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Profile from "../views/Profile.vue";
import Products_table from "../views/Products_table.vue";
import Cart from "../views/Cart.vue";
import Sells from "../views/Sells.vue";
import Buy from "../views/Buy.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/products',
      name: 'products',
      component: Products,
    },
    {
      path: '/product/:id',
      name: 'product',
      component: ProductDetail,
      props: true,
    },

    { path: "/login", component: Login },
    { path: "/register", component: Register },
    { path: "/profile", component: Profile },
    { path: "/", redirect: "/login" },
    { path: "/products_table", component: Products_table },
    { path: "/get_cart", component: Cart },
    { path: "/sells", component: Sells },
    { path: "/buy", component: Buy },
  ],
})

export default router
