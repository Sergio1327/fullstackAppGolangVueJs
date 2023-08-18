<template>
    <div>
        <stockModal />
        <stockListVue :stockList="stockList" />
    </div>
</template>

<script>
import stockListVue from '~/components/stockList.vue'
import stockModal from '~/components/stockModal.vue';

export default {
    data() {
        return {
            stockList: []
        }
    },
    components: {
        stockListVue,
        stockModal

    },
    methods: {
        async fetchStockList() {
            try {
                const response = await fetch('http://localhost:9000/stock_list');
                const responseData = await response.json();

                this.stockList = responseData.Data.stock_list

            } catch (error) {
                console.error('Error fetching warehouses:', error);
            }
        },
    },
    async mounted() {
        await this.fetchStockList();
        setInterval(this.fetchStockList, 5000)
    }
} 
</script>

