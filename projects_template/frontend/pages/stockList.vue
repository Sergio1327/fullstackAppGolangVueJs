<template>
    <div>
        <div class="is-flex is-align-items-center is-justify-content-flex-end">
            <button class="btn mb-5" @click="openModal" value="" type="submit">Добавить склад</button>
            <AddStockForm v-if="modalVisible" @closeModal="closeModal" :modalVisible="modalVisible" />
        </div>

        <b-table class="mt-6" :data="stockList" :hoverable="isHoverable" :striped="isStriped" :columns="columns"></b-table>
    </div>
</template>

<script>
import AddStockForm from '~/components/AddStockForm.vue';

export default {
    data() {
        return {
            modalVisible: false,
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
        AddStockForm

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
        openModal() {
            this.modalVisible = true
        },
        closeModal() {
            this.modalVisible = false
        }
    },
    async mounted() {
        await this.fetchStockList();
        setInterval(this.fetchStockList, 5000)
    }
} 
</script>

