<template>
    <div>
        <CardContent v-for="card in cards" :key="card.id" :card="card" />
    </div>
</template>

<script>

import CardContent from './ CardContent.vue';
export default {
    components: {
        CardContent
    },
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
}

</script>

