<template>
  <div class="container mt-5">
    <h1 class="text-center mb-4">Вход в систему</h1>
    <h5 class="text-center mb-4">После входа или выхода обновить страницу F5</h5>
    <div class="row justify-content-center">
      <div class="col-md-6">
        <form @submit.prevent="login" class="bg-light p-4 rounded shadow-sm">
          <div class="form-group">
            <label for="email">Email:</label>
            <input type="email" id="email" v-model="email" class="form-control" required />
          </div>
          <div class="form-group">
            <label for="password">Пароль:</label>
            <input type="password" id="password" v-model="password" class="form-control" required />
          </div>
          <button type="submit" class="btn btn-primary btn-block mt-3 button5">
            Войти
          </button>
          <button class="btn btn-primary btn-block mt-3">
            <router-link class="nav-link" to="/register">
              <span>Нет аккаунта?</span>
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
      email: "",
      password: "",
      errorMessage: "",
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post(
          "http://localhost:8080/login",
          {
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
        console.log("Ответ сервера:", response.data.data.username);

        const userData = {
          id: response.data.data.id,
          email: response.data.data.email,
          name: response.data.data.username,
          photo_url: response.data.data.photo_url,
          session: response.data.data.session,
        };


        console.log("Ответ сервера:", userData);
        localStorage.setItem("user", JSON.stringify(userData));
        console.log("Ответ сервера:", response.data);
        alert("Вы успешно вошли в систему!");

        localStorage.setItem("page_redir_path", "/profile");

        // this.$router.push("/profile");
        window.location.reload();
      } catch (error) {
        if (error.response) {
          console.error("Ошибка сервера:", error.response.data);
          this.errorMessage =
            error.response.data.message || "Ошибка входа. Проверьте данные.";
        } else {
          console.error("Ошибка сети:", error);
          this.errorMessage = "Ошибка сети. Попробуйте снова.";
        }
      }
    },
  },

  created() {
    const redirectPath = localStorage.getItem("page_redir_path");
    if (redirectPath) {
      localStorage.removeItem("page_redir_path");
      this.$router.push(redirectPath);
    }
  },
};
</script>

<style scoped>
.error {
  color: red;
}

.button5 {
  margin-right: 30px;
}
</style>
