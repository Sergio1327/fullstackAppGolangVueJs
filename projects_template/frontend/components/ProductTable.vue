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

                    <td><button @click="openDetailsModal(p.ProductID)" class="py-2 px-3">Просмотр деталей</button></td>
                    <td><button class="py-2 px-3">Добавить цены</button></td>
                    <td><button class="py-2 px-3">Добавить на склад</button></td>
                </tr>
            </tbody>

        </table>
        <ProductDetailsModal v-if="showModalDetails" :modalVisible="showModalDetails" :variantList="variantList"
            @closeModal="closeDetailsModal" />
    </div>
</template>

<script>
import ProductDetailsModal from './ProductDetailsModal.vue'
export default {
    components: {
        ProductDetailsModal
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
            resp: "",
            variantList: []
        }
    },
    methods: {
        closeDetailsModal() {
            this.showModalDetails = false
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
                console.log(responseData.Data.product_info.VariantList)
                this.variantList = responseData.Data.product_info.VariantList

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