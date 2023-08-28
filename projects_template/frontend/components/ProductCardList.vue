<script>
import productCard from './ProductCard.vue';
export default {


    components: {
        productCard
    },
    props: {
        productList: {
            type: Array,
            required: true,
        }
    }
}
</script>

<template>
    <div class="mt-6 ">
        <productCard v-for="p in productList" :key="p.ProductID" :productInfo="p" />
    </div>
</template> 