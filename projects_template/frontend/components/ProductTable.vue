<template>
    <div>
        <table border="2">
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

                    <td><button @click="openDetailsModal(p.ProductID)" class="py-2 button is-primary px-3">Просмотр
                            деталей</button></td>
                    <td><button @click="openPriceDetails(p.ProductID)" class="py-2 button is-warning px-3">Добавить
                            цены</button></td>
                    <td><button @click="OpenStockModal(p.ProductID)" class="py-2 button is-light px-3">Добавить на
                            склад</button></td>
                </tr>
            </tbody>

        </table>
        <ProductDetailsModal v-if="showModalDetails" :modalVisible="showModalDetails" :variantList="variantList"
            @closeModal="closeDetailsModal" />
        <ProductPriceModal v-if="showPriceModal" :modalVisible="showPriceModal" :options="VariantIDs"
            @closeModal="closePriceModal" />
        <ProductStockModal v-if="showStockModal" :modalVisible="showStockModal" :options="VariantIDs"
            @closeModal="closeStockModal" :storage-options="stockList" />
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
            VariantIDs: [],
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
            this.showModalDetails = true
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

            } catch (error) {
                console.error(error)
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
                this.VariantIDs = data.map(e => {
                    return {
                        Value: e.variant_id,
                        Option: e.variant_id
                    }
                })
                this.VariantIDs.sort((a, b) => a.Option - b.Option)
                this.showPriceModal = true

            } catch (error) {
                console.error(error)
            }
        },
        async OpenStockModal(ProductID) {
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
                this.VariantIDs = data.map(e => {
                    return {
                        Value: e.variant_id,
                        Option: e.variant_id
                    }
                })

                this.VariantIDs.sort((a, b) => a.Option - b.Option)

                const response2 = await fetch("http://127.0.0.1:9000/stock_list", {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    },

                })
                const responseData2 = await response2.json()
                const data2 = responseData2.Data.stock_list

                this.stockList = data2.map(e => {
                    return {
                        Option: e.StorageID,
                        Value: e.StorageID
                    }
                })
                this.showStockModal = true
            } catch (error) {
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