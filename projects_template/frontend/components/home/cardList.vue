<template>
    <div class="">
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
                    imgLink: "@assets/productList.jpg"
                },
                {
                    id: 2,
                    title: "Склады",
                    descr: "описание страницы складов",
                    imgLink: "@assets/stockList.jpg"
                },
                {
                    id: 3,
                    title: "Продажи",
                    descr: "описание страницы продаж",
                    imgLink: "@assets/saleList.png"
                }
            ]
        }
    },
}

</script>