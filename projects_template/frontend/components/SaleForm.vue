<template>
  <div>
    <div class="field">
      <label class="label">Дата начала продаж</label>
      <div class="control">
        <input class="input" type="datetime-local" v-model="startDate">
      </div>
    </div>

    <div class="field">
      <label class="label">Дата конца продаж</label>
      <div class="control">
        <input class="input" type="datetime-local" v-model="endDate">
      </div>
    </div>

    <div class="field">
      <label class="label">Лимит вывода</label>
      <div class="control">
        <div class="select">
          <select v-model="limit">
            <option value="1">1</option>
            <option value="5">5</option>
            <option value="10">10</option>
            <option value="50">50</option>
          </select>
        </div>
      </div>
    </div>

    <div class="field">
      <label class="label">Название продукта</label>
      <div class="control">
        <input class="input" type="text" v-model="productName">
      </div>
    </div>

    <div class="field">
      <label class="label">ID склада</label>
      <div class="control">
        <input class="input" type="number" v-model="storageId">
      </div>
    </div>

    <div class="field is-grouped">
      <div class="control">
        <button class="button is-link" @click="sendRequest">Найти продажи</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      startDate: null,
      endDate: null,
      limit: 1,
      productName: "",
      storageId: 1,

    };
  },
  methods: {
    async sendRequest() {

      const formattedStartDate = this.formatDate(this.startDate);
      const formattedEndDate = this.formatDate(this.endDate);

      const requestData = {
        start_date: formattedStartDate,
        end_date: formattedEndDate,
        limit: parseInt(this.limit),
        product_name: this.productName,
        storage_id: parseInt(this.storageId)
      };
      console.log(requestData)
      try {
        const response = await fetch("http://127.0.0.1:9000/sales", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(requestData),
        });

        const responseData = await response.json();
        const data = responseData.Data.sale_list
        console.log(data)
        this.$emit("data", data)

      } catch (error) {
        console.error("Ошибка при отправке запроса:", error);
        console.log(error)
      }
    },
    formatDate(dateTime) {
      const date = new Date(dateTime);
      const year = date.getUTCFullYear();
      const month = String(date.getUTCMonth() + 1).padStart(2, "0");
      const day = String(date.getUTCDate()).padStart(2, "0");
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getUTCMinutes()).padStart(2, "0");
      const seconds = String(date.getUTCSeconds()).padStart(2, "0");

      return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}+05:00`;
    },
  },
};
</script>