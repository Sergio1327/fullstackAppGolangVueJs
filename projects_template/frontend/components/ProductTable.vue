<template>
    <div class="table-container">
        <div class="table  is-striped is-narrow is-hoverable is-fullwidth">
            <table border="0">
                <thead>
                    <tr>
                        <th>ID продукта</th>
                        <th>Название продукта</th>
                        <th>Описание продукта</th>
                        <th>Детали</th>
                        <th>Цены</th>
                        <th>Добавить на склад</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="p in productList" :key="p.ProductID">
                        <td>{{ p.ProductID }}</td>
                        <td>{{ p.Name }}</td>
                        <td>{{ p.Descr }}</td>

                        <td><button @click="openDetailsModal(p.ProductID)" class=" button is-primary ">Просмотр
                                деталей</button></td>
                        <td><button @click="openPriceDetails(p.ProductID)" class=" button is-warning ">Добавить
                                цены</button></td>
                        <td><button @click="OpenStockModal(p.ProductID)" class=" button is-light ">Добавить на
                                склад</button></td>
                    </tr>
                </tbody>

            </table>

            <ProductDetailsModal v-if="showModalDetails" :modalVisible="showModalDetails" :variantList="variantList"
                @closeModal="closeDetailsModal" />
            <ProductPriceModal v-if="showPriceModal" :modalVisible="showPriceModal" :options="variantIDs"
                @closeModal="closePriceModal" />
            <ProductStockModal v-if="showStockModal" :modalVisible="showStockModal" :options="variantIDs"
                @closeModal="closeStockModal" :storage-options="stockList" />
        </div>
    </div>
</template>

<script>
import ProductDetailsModal from './ProductDetailsModal.vue'
import ProductPriceModal from './ProductPriceModal.vue'
import ProductStockModal from './ProductStockModal.vue'

export default {
    components: {
        ProductDetailsModal,
        ProductPriceModal,
        ProductStockModal
    },

    props: {
        productList: {
            type: Array,
            required: true
        }
    },

    data() {
        return {
            showModalDetails: false,
            showPriceModal: false,
            showStockModal: false,

            resp: "",

            variantList: [],
            variantIDs: [],
            stockList: []
        }
    },

    methods: {
        closeDetailsModal() {
            this.showModalDetails = false
        },

        closePriceModal() {
            this.showPriceModal = false
        },

        closeStockModal() {
            this.showStockModal = false
        },

        async openDetailsModal(ProductID) {
            const url = `http://127.0.0.1:9000/product/${ProductID}`

            try {
                const response = await fetch(url, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                })

                const responseData = await response.json()
                this.variantList = responseData.Data.product_info.VariantList
                this.showModalDetails = true
            }
            catch (error) {
                console.error(error)
                this.$buefy.snackbar.open(error)
            }
        },

        async openPriceDetails(ProductID) {
            const url = `http://127.0.0.1:9000/product/${ProductID}`
            try {
                const response = await fetch(url, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                })
                const responseData = await response.json()

                const data = responseData.Data.product_info.VariantList
                this.variantIDs = data.map(e => {
                    return {
                        Value: e.variant_id,
                        Option: e.variant_id,
                        weight: e.weight,
                        unit: e.unit
                    }
                })
                this.variantIDs.sort((a, b) => a.Option - b.Option)
                this.showPriceModal = true

            }
            catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        },

        async OpenStockModal(ProductID) {
            const url = `http://127.0.0.1:9000/product/${ProductID}`

            try {
                const variantResponse = await fetch(url, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                })
                const variantResponseData = await variantResponse.json()

                const variantData = variantResponseData.Data.product_info.VariantList
                this.variantIDs = variantData.map(e => {
                    return {
                        Value: e.variant_id,
                        Option: e.variant_id,
                        weight: e.weight,
                        unit: e.unit
                    }
                })

                this.variantIDs.sort((a, b) => a.Option - b.Option)

                const stockResponse = await fetch("http://127.0.0.1:9000/stock_list", {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    },

                })
                const stockResponseData = await stockResponse.json()
                const stockData = stockResponseData.Data.stock_list

                this.stockList = stockData.map(e => {
                    return {
                        StorageName: e.StorageName,
                        Option: e.StorageID,
                        Value: e.StorageID,
                    }
                })

                this.showStockModal = true
            }
            catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        }

    }
}

</script>

<style scoped>
th,
td {
    text-align: center !important;
    padding: 1% 2%;

}
table {
    margin-top: 50px;
    width: 100%;
}

</style>