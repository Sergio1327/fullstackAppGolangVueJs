<template>
  <div class="card ">
    <div class="card-content is-flex is-align-items-center is-justify-content-space-between ">
      <div class="content column is-text-center">
        <div>ID продажи : <strong> {{ data.SaleID }}</strong></div>
        <div>Название продукта : <strong> {{ data.ProductName }}</strong></div>
        <div>ID Варианта продукта : <strong> {{ data.variant_id }}</strong></div>
        <div>ID склада : <strong> {{ data.storage_id }}</strong></div>
        <div>Продано в : <strong> {{ data.SoldAt }}</strong></div>
        <div>Кол-во : <strong> {{ data.quantity }}</strong></div>
        <div>Общая цена : <strong> {{ data.TotalPrice }}</strong></div>
      </div>
      <div class="image is-128x128">
        <img src="@/assets/saleList.png" class="is-img" alt="">
      </div>
    </div>
  </div>
</template>
  
<script>
export default {
  props: {
    data: {
      type: Object,
      required: true
    }
  }
};
</script>
  
<style scoped>
.image {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card {
  margin-bottom: 20px;
  border-radius: 5px;
  box-shadow: 1px 1px 8px orange;
}

.card-content {
  padding: 1rem;
}

.content>div>strong {
  font-weight: 600;
  color: #000;
}

.title.is-4 {
  margin-bottom: 5px;
}

.subtitle.is-6 {
  color: #999;
}
</style>