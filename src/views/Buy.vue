<template>
  <div>
    <h1 style="display: flex; justify-content: center; margin: 30px;">Мои покупочки</h1>

    <div v-if="loading">Загрузка...</div>
    <div v-if="error" class="error">{{ error }}</div>

    <div class="filters mb-3">
      <div class="row">
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_id" placeholder="Фильтр по ID продажи" class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_u_id" placeholder="Фильтр по ID пользователя" class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_un_id" placeholder="Фильтр по имени пользователя"
            class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_p_id" placeholder="Фильтр по ID товара" class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_pn_id" placeholder="Фильтр по названию товара"
            class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_q_id" placeholder="Фильтр по количеству" class="form-control" />
        </div>
        <div class="col-md-3 mb-2">
          <input type="text" v-model="filters.sell_tp_id" placeholder="Фильтр по цене" class="form-control" />
        </div>
        <div class="col-md-6 mb-2">
          <label>Фильтр по дате продажи:</label>
          <div class="d-flex">
            <input type="date" v-model="filters.startDate" class="form-control me-2" />
            <input type="date" v-model="filters.endDate" class="form-control" />
          </div>
        </div>
        <!-- <div class="col-md-3 mb-2">
          <label>Фильтр по номеру корзины:</label>
          <input type="text" v-model="filters.CartID" placeholder="Фильтр по CartID" class="form-control" />
        </div> -->
      </div>
    </div>
    <button class="btn btn-primary mt-3" @click="printOrders1(filteredOrders, totalFilteredAmount)">Сгенерировать PDF</button>

    <div class="table-responsive">
      <table v-if="filteredOrders.length > 0" class="table table-striped table-bordered table-hover">
        <thead class="thead-dark">
          <tr>
            <th>ID продажи</th>
            <th>ID пользователя</th>
            <th>Имя пользователя</th>
            <th>ID товара</th>
            <th>Название товара</th>
            <th>Количество</th>
            <th>Цена</th>
            <th>Дата продажи</th>
            <!-- <th>CartID</th> -->
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in filteredOrders" :key="order.id">
            <td>{{ order.sell_id }}</td>
            <td>{{ order.sell_u_id }}</td>
            <td>{{ order.sell_un_id }}</td>
            <td>{{ order.sell_p_id }}</td>
            <td>{{ order.sell_pn_id }}</td>
            <td>{{ order.sell_q_id }}</td>
            <td>{{ formatNumber(order.sell_tp_id) }} ₽</td>
            <td>{{ order.autoCreateTime }}</td>
            <!-- <td>{{ order.CartID }}</td> -->
          </tr>
          <tr>
            <td>Итого: </td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td>{{ formatNumber(totalFilteredAmount) }} ₽</td>
            <td></td>
            <!-- <td></td> -->
          </tr>
        </tbody>
      </table>
      <p v-else>Нет данных, соответствующих фильтру.</p>
    </div>
  </div>
</template>

<script>

import 'jspdf-autotable';
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
      filters: {
        sell_id: "",
        sell_u_id: "",
        sell_un_id: "",
        sell_p_id: "",
        sell_pn_id: "",
        sell_q_id: "",
        sell_tp_id: "",
        startDate: "",
        endDate: "",
        CartID: ""
      },
      orders: [],
      loading: true,
      error: null,
    };
  },

  computed: {
    printOrders() {
      const sellData = this.orders.map(order => ({
        sell_id: order.sell_id,
        sell_u_id: order.sell_u_id,
        user_name: order.sell_un_id,
        product_id: order.sell_p_id,
        product_name: order.sell_pn_id,
        quantity: order.sell_q_id,
        total_price: order.sell_tp_id,
        sale_date: order.autoCreateTime,
        cart_id: order.CartID
      }));

      console.log(sellData)
      fetch('http://localhost:8080/generate-pdf-sell', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sellData: sellData })
      })
        .then(response => {
          if (response.ok) {
            return response.blob();
          }
          throw new Error("Failed to generate PDF");
        })
        .then(blob => {
          const link = document.createElement('a');
          link.href = URL.createObjectURL(blob);
          link.download = 'orders.pdf';
          link.click();
        })
        .catch(error => console.error('Error:', error));
    },

    filteredOrders() {
      return this.orders.filter(order => {
        // Фильтрация по ID продажи
        const sellIdMatch = this.filters.sell_id
          ? order.sell_id.toString().includes(this.filters.sell_id)
          : true;

        // Фильтрация по ID пользователя
        const userIdMatch = this.filters.sell_u_id
          ? order.sell_u_id.toString().includes(this.filters.sell_u_id)
          : true;

        // Фильтрация по имени пользователя
        const userNameMatch = this.filters.sell_un_id
          ? order.sell_un_id.toLowerCase().includes(this.filters.sell_un_id.toLowerCase())
          : true;

        // Фильтрация по ID товара
        const productIdMatch = this.filters.sell_p_id
          ? order.sell_p_id.toString().includes(this.filters.sell_p_id)
          : true;

        // Фильтрация по названию товара
        const productNameMatch = this.filters.sell_pn_id
          ? order.sell_pn_id.toLowerCase().includes(this.filters.sell_pn_id.toLowerCase())
          : true;

        // Фильтрация по количеству
        const quantityMatch = this.filters.sell_q_id
          ? order.sell_q_id.toString().includes(this.filters.sell_q_id)
          : true;

        // Фильтрация по цене
        const priceMatch = this.filters.sell_tp_id
          ? order.sell_tp_id.toString().includes(this.filters.sell_tp_id)
          : true;

        // Фильтрация по диапазону дат
        const startDate = this.filters.startDate ? new Date(this.filters.startDate) : null;
        const endDate = this.filters.endDate ? new Date(this.filters.endDate) : null;
        const orderDate = new Date(order.autoCreateTime);

        const dateMatch =
          (!startDate || orderDate >= startDate) &&
          (!endDate || orderDate <= endDate);

        // Фильтрация по CartID
        const cartIdMatch = this.filters.CartID
          ? order.CartID.toLowerCase().includes(this.filters.CartID.toLowerCase())
          : true;

        // Возвращаем заказ, если все условия фильтрации совпадают
        return (
          sellIdMatch &&
          userIdMatch &&
          userNameMatch &&
          productIdMatch &&
          productNameMatch &&
          quantityMatch &&
          priceMatch &&
          dateMatch &&
          cartIdMatch
        );
      });
    },
    // Общая сумма для отфильтрованных заказов
    totalFilteredAmount() {
      return this.filteredOrders.reduce((sum, order) => sum + order.sell_tp_id, 0);
    },
    totalAmount() {
      // Считаем сумму всех товаров
      console.log(this.orders.reduce((sum, order) => sum + (order.sell_q_id * parseInt(order.sell_tp_id)), 0))

      return this.orders.reduce((sum, order) => {
        const value = Number(order.sell_tp_id); // Преобразуем значение в число
        return sum + (isNaN(value) ? 0 : value); // Если значение не число, добавляем 0
      }, 0);
    }
  },
  mounted() {
    this.fetchOrders();
  },
  methods: {

    printOrders1(filteredOrders, totalFilteredAmount) {
      const printContent = `
  <html>
    <head>
      <title>Отчет по заказам</title>
      <style>
        /* Настройка страницы для печати */
        @page {
          size: A4 landscape; /* Горизонтальная ориентация */
          margin: 20mm;
        }

        body {
          font-family: 'Arial', sans-serif;
          margin: 0;
          padding: 20px;
          background-color: #f9f9f9;
          color: #333;
          font-size: 14px;
        }

        h1 {
          text-align: center;
          font-size: 28px;
          color: #4a90e2;
          margin-bottom: 20px;
        }

        table {
          width: 100%;
          border-collapse: collapse;
          margin-top: 20px;
          border: 1px solid #ddd;
          background-color: #fff;
        }

        th, td {
          padding: 12px 15px;
          text-align: center;
          border: 1px solid #ddd;
          font-size: 14px;
        }

        th {
          background-color: #3498db;
          color: white;
          font-weight: bold;
        }

        td {
          background-color: #f8f8f8;
        }

        .total-row {
          font-weight: bold;
          background-color: #e6f1ff;
        }

        .total-cell {
          font-size: 16px;
          font-weight: bold;
          color: #2c3e50;
        }

        .footer {
          font-size: 12px;
          color: #888;
          text-align: center;
          margin-top: 30px;
        }

        /* Печать */
        @media print {
          body {
            background-color: #ffffff;
            margin: 0;
            padding: 0;
          }

          table {
            width: 100%;
            border-collapse: collapse;
          }

          th, td {
            padding: 10px;
            text-align: center;
          }

          .footer {
            font-size: 10px;
          }
        }
      </style>
    </head>
    <body>
      <h1>Отчет по заказам</h1>
      
      <table>
        <thead>
          <tr>
            <th>ID продажи</th>
            <th>ID пользователя</th>
            <th>Имя пользователя</th>
            <th>ID товара</th>
            <th>Название товара</th>
            <th>Количество</th>
            <th>Цена</th>
            <th>Дата продажи</th>
          </tr>
        </thead>
        <tbody>
          <!-- Итерация по filteredOrders -->
          ${this.filteredOrders.map(order => `
            <tr>
              <td>${order.sell_id}</td>
              <td>${order.sell_u_id}</td>
              <td>${order.sell_un_id}</td>
              <td>${order.sell_p_id}</td>
              <td>${order.sell_pn_id}</td>
              <td>${order.sell_q_id}</td>
              <td>${this.formatNumber(order.sell_tp_id)} ₽</td>
              <td>${order.autoCreateTime}</td>
            </tr>
          `).join('')}
          <!-- Итоговая строка -->
          <tr class="total-row">
            <td colspan="6" class="total-cell">Итого:</td>
            <td class="total-cell">${this.formatNumber(totalFilteredAmount)} ₽</td>
            <td></td>
          </tr>
        </tbody>
      </table>

      <div class="footer">
        <p>© 2025 TechStore Рожков Григорий Олегович. Все права не защищены((((.</p>
      </div>
    </body>
  </html>
  `;
      const printWindow = window.open('', '_blank', 'width=800,height=600');
      printWindow.document.write(printContent);
      printWindow.document.close();
      printWindow.print();
    },

    formatNumber(value) {
      return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
    },

    async fetchOrders() {
      try {
        const userID = JSON.parse(localStorage.getItem('user')).id;
        console.log("ЭТО ИД", userID)

        if (!userID) {
          console.error('userID не найден в localStorage');
          this.error = 'Пользователь не авторизован.';
          return;
        }
        const response = await axios.get(`http://localhost:8080/buy/${userID}`);

        this.orders = response.data;
        console.log(this.orders)
      } catch (err) {
        this.error = "Ошибка при загрузке заказов: " + err.message;
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.table {
  margin-top: 20px;
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

.error {
  color: red;
}
</style>