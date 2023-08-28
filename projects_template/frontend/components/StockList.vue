<template>
    <div>
        <ul>
            <stockItem class="mb-6" v-for="stock in stockList" :key="stock.StorageID" :stock="stock" />
        </ul>
    </div>
</template>

<script>
import stockItem from './StockItem.vue';

export default {
    components: {
        stockItem
    },
    props: {
        stockList: {
            type: Array,
            required: true
        },
    }
}
</script>