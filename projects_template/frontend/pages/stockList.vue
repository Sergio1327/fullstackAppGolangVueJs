<template>
    <div>
        <stockModal />
        <b-table class="mt-6" :data="stockList" :hoverable="isHoverable" :striped="isStriped" :columns="columns"></b-table>
    </div>
</template>

<script>
import stockModal from '~/components/StockModal.vue';

export default {
    data() {
        return {
            stockList: [],
            columns: [
                {
                    field: "StorageID",
                    label: "ID склада",
                    centered: true,
                    width: 200,
                    height: 300
                },
                {
                    field: "StorageName",
                    label: "Название склада",
                    centered: true
                },
            ],
            isBordered: true,
            isHoverable: true,
            isStriped: true
        }

    },
    components: {
        stockModal

    },
    methods: {
        async fetchStockList() {
            try {
                const response = await fetch('http://localhost:9000/stock_list');
                const responseData = await response.json();

                this.stockList = responseData.Data.stock_list
                console.log(this.stockList)
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

