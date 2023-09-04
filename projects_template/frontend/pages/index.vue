<template>
  <section>
    <CardContent class="mt-6" v-for="card in cards" :key="card.id" :card="card" />
  </section>
</template>

<script>
import CardContent from '~/components/ CardContent.vue';

export default {
  data() {
    return {
      cards: [
        {
          id: 1,
          title: "Продукты",
          descr: "описание страницы продуктов",
          imgLink: require("@/assets/productList.jpg"),
          link: "/productList"
        },
        {
          id: 2,
          title: "Склады",
          descr: "описание страницы складов",
          imgLink: require("@/assets/stockList.jpg"),
          link: "/stockList"
        },
        {
          id: 3,
          title: "Продажи",
          descr: "описание страницы продаж",
          imgLink: require("@/assets/saleList.png"),
          link: "/saleList"
        }
      ]
    }
  },
  components: {
    CardContent,
  },
}
</script>


<style></style>