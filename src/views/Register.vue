<template>
  <div class="container mt-5">
    <h1 class="text-center">Регистрация</h1>
    <div class="row justify-content-center">
      <div class="col-md-6">
        <form @submit.prevent="register" class="bg-light p-4 rounded shadow-sm">
          <div class="form-group">
            <label for="username">Имя пользователя:</label>
            <input type="text" id="username" v-model="username" class="form-control" required />
          </div>
          <div class="form-group">
            <label for="email">Email:</label>
            <input type="email" id="email" v-model="email" class="form-control" required />
          </div>
          <div class="form-group passwor">
            <label for="password">Пароль:</label>
            <input type="password" id="password" v-model="password" class="form-control" required />
          </div>
          <button type="submit" class="btn btn-primary btn-block button5">
            Зарегистрироваться
          </button>
          <button class="btn btn-primary btn-block">
            <router-link class="nav-link" to="/login">
              <span>У меня есть аккаунт</span>
            </router-link>
          </button>
          <p v-if="errorMessage" class="text-danger mt-3">{{ errorMessage }}</p>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      username: "",
      email: "",
      password: "",
      errorMessage: "",
    };
  },
  methods: {
    async register() {
      try {
        const response = await axios.post(
          "http://localhost:8080/register",
          {
            username: this.username,
            email: this.email,
            password: this.password,
          },
          {
            headers: {
              "Content-Type": "application/json",
            },
            withCredentials: true,
          }
        );

        console.log("Ответ сервера:", response.data);
        alert("Регистрация успешна!");

        this.$router.push("/login").then(() => {
        location.reload();
        });
      } catch (error) {
        if (error.response) {
          console.error("Ошибка сервера:", error.response.data);
          this.errorMessage =
            error.response.data.message || "Ошибка регистрации";
        } else {
          console.error("Ошибка сети:", error);
          this.errorMessage = "Ошибка сети. Попробуйте снова.";
        }
      }
    },
  },
};
</script>

<style scoped>
.container {
  margin-top: 50px;
}

.passwor {
  margin-bottom: 30px;
}

.button5 {
  margin-right: 30px;
}
</style>
