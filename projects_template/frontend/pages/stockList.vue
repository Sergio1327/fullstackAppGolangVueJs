<template>
    <section>
        <div class="mt-6 is-flex is-align-items-center is-justify-content-flex-end">
            <button class="button py-3 px-5 is-warning mb-5" @click="openModal" value="" type="submit">Добавить
                склад</button>
            <AddStockForm v-if="modalVisible" @fetchStockList="fetchStockList" @closeModal="closeModal"
                :modalVisible="modalVisible" />
        </div>

        <b-table class="mt-6 table" :data="stockList" :hoverable="isHoverable" :striped="isStriped"
            :columns="columns"></b-table>
    </section>
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
                const response = await fetch('http://localhost:9000/stock_list', {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                const responseData = await response.json();

                this.stockList = responseData.Data.stock_list
                
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
    }
} 
</script>


<style scoped >
.table {
    margin-top: 100px !important;
    margin-bottom: 350px !important;
}
</style>