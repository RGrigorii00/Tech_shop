<template>

  <link href="https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.css" rel="stylesheet">

  <div>
    <h1 style="display: flex; justify-content: center; margin: 30px;">Корзина</h1>

    <div v-if="cart.length == 0">
      <div style="display: flex; justify-content: center; margin: 30px;">Корзина пуста, пока что здесь нет товаров :(
      </div>
    </div>
    <div v-else>
      <div  class="table-responsive">
      <table class="table table-striped">
        <thead>
          <tr>
            <th>Номер товара</th>
            <th>Название товара</th>
            <th>Цена</th>
            <th>Количество</th>
            <th>Цена за все</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in cart" :key="item.id">
            <!-- <td>{{ item }}</td> -->
            <td>{{ item.product_id }}</td>
            <td>{{ item.product_name }}</td>
            <td>{{ formatNumber(item.product_price) }} ₽</td>

            <td>
              <input type="number" v-model="item.quantity" :min="1" class="form-control" />
            </td>
            <td>{{ formatNumber(item.product_price * item.quantity) }} ₽</td>
            <td>
              <button @click="removeFromCart(item.cart_id)" class="btn btn-danger btn-sm">Удалить</button>
            </td>
          </tr>
        </tbody>
      </table>
      </div>

      <button style="margin-right: 10px;" @click="printCart" class="btn btn-primary">
        Скачать PDF с корзиной (с сервера)
      </button>

      <button @click="printTable" class="btn btn-primary">
        Скачать PDF с корзиной (с клиента)
      </button>

      <div class="d-flex justify-content-between mt-4">
        <h4>Итого: {{ formatNumber(total) }} ₽</h4>

        <button @click="makePurchase" class="btn btn-success">Купить</button>

        <button @click="clearCart" class="btn btn-warning">Очистить корзину</button>
      </div>
    </div>
  </div>
</template>

<script>
import toastr from "toastr";
import jschardet from 'jschardet';
// Функция для конвертации текста из неизвестной кодировки в UTF-8
function convertText(text) {
  var detected = jschardet.detect(text);
  console.log("Detected encoding: ", detected.encoding);

  if (detected.encoding && detected.encoding !== 'UTF-8') {
    try {
      const decoder = new TextDecoder(detected.encoding, { fatal: true });
      const encoded = new TextEncoder().encode(text);
      const decoded = decoder.decode(encoded);
      return decoded;
    } catch (e) {
      console.error("Error decoding text:", e);
      return "Ошибка при декодировании текста.";
    }
  }
  return text;
}
import axios from "axios";

export default {
  data() {
    return {
      cart: [],
      loading: false,
      itog: null,
    };
  },
  computed: {
    total() {
      return this.cart.reduce((sum, item) => sum + item.product_price * item.quantity, 0);
    },
  },
  methods: {

    printTable() {
    // Получаем HTML таблицы
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
          .total-row { font-weight: bold; background-color: #e6f1ff; }
          .total-cell { font-weight: bold; color: #2c3e50; }
        </style>
      </head>
      <body>
        <div class="table-responsive">
          <table class="table table-striped">
            <thead>
              <tr>
                <th>Номер товара</th>
                <th>Название товара</th>
                <th>Цена</th>
                <th>Количество</th>
                <th>Цена за все</th>
              </tr>
            </thead>
            <tbody>
              ${this.cart.map(item => `
                <tr>
                  <td>${item.product_id}</td>
                  <td>${item.product_name}</td>
                  <td>${this.formatNumber(item.product_price)} ₽</td>
                  <td>${item.quantity}</td>
                  <td>${this.formatNumber(item.product_price * item.quantity)} ₽</td>
                </tr>
              `).join('')}
            </tbody>
            <tfoot>
              <tr class="total-row">
                <td colspan="4" class="total-cell">Итого: </td>
                <td colspan="2" class="total-cell">${this.formatNumber(this.cart.reduce((total, item) => total + item.product_price * item.quantity, 0))} ₽</td>
              </tr>
            </tfoot>
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
    itog() {
      itog = this.total();
    },

    printCart() {
      const cartData = this.cart.map(item => ({
        cart_id: convertText(item.cart_id),
        user_name: convertText(item.user_name),
        product_name: convertText(item.product_name),
        product_price: convertText(item.product_price),
        quantity: convertText(item.quantity)
      }));

      fetch('http://localhost:8080/generate-pdf', {
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
    },

    async makePurchase() {
      try {
        const userId = localStorage.getItem("user");
        const user = JSON.parse(userId);

        if (!user || !user.id) {
          console.error("Ошибка: user_id не найден в localStorage");
          alert("Ошибка: не удалось найти информацию о пользователе.");
          return;
        }

        // Проверяем, что корзина не пуста
        if (!this.cart || this.cart.length === 0) {
          console.error("Ошибка: корзина пуста");
          alert("Корзина пуста. Пожалуйста, добавьте товары в корзину.");
          return;
        }

        const cartItems = this.cart.map(item => ({
          cart_id: item.cart_id,
          product_id: item.product_id,
          quantity: item.quantity,
          total_price: item.product_price * item.quantity
        }));
        console.log("User ID:", user.id);
        console.log("Cart Items:", cartItems);

        const response = await axios.post("http://localhost:8080/make_purchase", {
          user_id: user.id,
          cart_items: cartItems,
        });

        if (response.status === 200) {
          toastr.success("Покупка совершена!", "Успех", {
            closeButton: true,
            progressBar: true,
            positionClass: "toast-bottom-right",
          });

          this.clearCart();
          console.log("Покупка успешно завершена:", response.data);
        } else {
          console.error("Ошибка при оформлении покупки:", response.data);
          alert("Произошла ошибка при оформлении покупки.");
        }

      } catch (error) {

        toastr.error("Не удалось совершить покупку.", error, {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
        console.error("Ошибка при отправке запроса:", error);
        alert("Произошла ошибка. Попробуйте снова.");
      }
    },

    formatNumber(value) {
      return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
    },
    async fetchCart() {
      this.loading = true;
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
        const response = await axios.get("http://localhost:8080/get_cart", {
          params: { user_id: z },
        });
        console.log("Респонсе", response.data)
        console.log(this.item)
        this.cart = response.data;
      } catch (error) {
        console.error("Ошибка при загрузке корзины:", error);
      } finally {
        this.loading = false;
      }
    },

    async removeFromCart(cart_id) {
      try {
        const url = `http://localhost:8080/remove_from_cart?id=${cart_id}`;

        await axios.delete(url);
        toastr.success("Товар удален из корзины!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });

        this.fetchCart();
        console.log("МАЯ КАРЗИНААА", this.cart)
        console.log("МАЯ КАРЗИНААА", this.cart.length)
        if (this.cart && this.cart.length == 1) { window.location.reload(); }
      } catch (error) {
        console.error("Ошибка при удалении товаров:", error);

        toastr.error("Не удалось удалить товар из корзины.", error, {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
      }
    },

    async clearCart() {
      try {
        const userId = localStorage.getItem("user");
        const user = JSON.parse(userId);
        console.log("Парсинг user:", user);
        console.log(user.id)

        if (!userId || !user.id) {
          console.error("Ошибка: user_id не найден в localStorage");
          return;
        }

        await axios.delete("http://localhost:8080/clear_cart", {
          params: { user_id: user.id },
        });
        toastr.success("Корзина очищена!", "Успех", {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });

        window.location.reload();
        this.fetchCart();
      } catch (error) {
        toastr.error("Не удалось очистить корзину.", error, {
          closeButton: true,
          progressBar: true,
          positionClass: "toast-bottom-right",
        });
        console.error("Ошибка при очистке корзины:", error);
      }
    },
  },
  created() {
    this.fetchCart();
    console.log(this.cart)
  },
};
</script>

<style scoped>
table {
  margin-top: 20px;
}
</style>