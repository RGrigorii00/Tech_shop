<template>
  <div class="container">
    <div class="row justify-content-center">
      <div>
        <div class="card shadow-sm">
          <div class="card-body">
            <h2 class="text-center">Профиль пользователя</h2>

            <div>
              <div style="margin-bottom: 30px;" class="text-center">
                <img :src="profile.user.photo_url || 'https://via.placeholder.com/150'" alt="Фото пользователя"
                  class="rounded-circle img-fluid" style="width: 150px; height: 150px; object-fit: cover;">
              </div>

              <div class="mb-4 upload-box" @dragover.prevent="onDragOver" @dragleave.prevent="onDragLeave"
                @drop.prevent="onDrop" :class="{ 'drag-active': isDragActive }">
                <h4>Обновить фото</h4>
                <input ref="fileInput" type="file" accept="image/*" @change="handleFileUpload"
                  class="form-control d-none" />
                <p>Перетащите файл сюда или <span @click="triggerFileInput" class="upload-link">выберите файл</span>.
                </p>
              </div>

              <div class="mb-4">
                <h4>Основная информация</h4>
                <ul class="list-group">
                  <li class="list-group-item"><strong>Баланс: </strong> {{ formatNumber(profile.user.balance) }} ₽</li>
                  <li class="list-group-item"><strong>Имя пользователя:</strong> {{ profile.user.username }}</li>
                  <li class="list-group-item"><strong>Email:</strong> {{ profile.user.email }}</li>
                  <li class="list-group-item"><strong>Роль:</strong> {{ profile.user.role }}</li>
                  <li class="list-group-item"><strong>Пол:</strong> {{ profile.user.gender }}</li>
                  <li class="list-group-item"><strong>Адрес:</strong> {{ profile.user.address }}</li>
                  <li class="list-group-item"><strong>Фотки:</strong> {{ profile.user.photo_url }}</li>
                </ul>
              </div>

              <div class="mb-4">
                <h4>Сессия</h4>
                <ul class="list-group">
                  <li class="list-group-item"><strong>Ключ сессии:</strong> {{ profile.session.session_id }}</li>
                  <li class="list-group-item"><strong>Дата создания:</strong> {{ profile.session.created_at }}</li>
                </ul>
              </div>

              <div class="d-flex justify-content-center">
                <button @click="logout" class="btn btn-danger">Выход</button>
              </div>
            </div>
            <div>
              <p v-if="loading == false">Загрузка данных профиля...</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- <div v-if="this.profile.user.role == 'admin'">

    <button @click="downloadFile" class="btn btn-primary">Скачать архив (.rar)</button>
  </div> -->
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      profile: null,
      loading: false,
      errorMessage: "",
      isDragActive: false,
    };
  },
  created() {
    this.fetchProfile();
  },
  methods: {
    downloadFile() {
      axios({
        url: 'http://localhost:8080/daf',
        method: 'GET',
        responseType: 'blob',
      })
        .then((response) => {
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'your-archive.rar');
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
          window.URL.revokeObjectURL(url);
        })
        .catch((error) => {
          console.error('Ошибка при скачивании архива:', error);
        });
    },
    formatNumber(value) {
      return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
    },
    async fetchProfile() {
      try {
        const response = await axios.get("http://localhost:8080/profile", {
          withCredentials: true,
        });
        this.profile = response.data;
        this.loading = true;
        console.log(this.profile)
      } catch (error) {
        this.errorMessage =
          error.response?.data?.message || "Ошибка загрузки профиля";
      }
    },
    async handleFileUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.uploadFile(file);
      }
    },
    async uploadFile(file) {
      const formData = new FormData();
      formData.append("photo", file);
      formData.append("user_id", this.profile.user.id);

      try {
        console.log("Отправка данных:", ...formData.entries());
        const response = await axios.post(
          "http://localhost:8080/upload-photo",
          formData,
          {
            withCredentials: true,
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        this.profile.user.photo_url = response.data.photo_url;
        window.location.reload();
      } catch (error) {
        console.error("Ошибка запроса:", error.response || error.message);
        this.errorMessage =
          error.response?.data?.message || "Ошибка загрузки фото";
      }
    },
    onDragOver() {
      this.isDragActive = true;
    },
    onDragLeave() {
      this.isDragActive = false;
    },
    onDrop(event) {
      this.isDragActive = false;
      const file = event.dataTransfer.files[0];
      if (file) {
        this.uploadFile(file);
      }
    },
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    async logout() {
      try {
        await axios.get("http://localhost:8080/logout", {
          withCredentials: true,
        });
        this.profile = null;
        localStorage.removeItem('user');
        this.$router.push("/login").then(() => {
        location.reload();
        });
      } catch (error) {
        this.errorMessage =
          error.response?.data?.message || "Ошибка при выходе из системы";
      }
    },
  },
};
</script>

<style scoped>
.container {
  margin-top: 50px;
}

.card {
  border-radius: 8px;
}

.card-body {
  padding: 30px;
}

.list-group-item {
  font-size: 1rem;
  border: none;
}

.list-group-item strong {
  color: #333;
}

img {
  object-fit: cover;
}

button {
  width: 100%;
  padding: 12px;
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
</style>
