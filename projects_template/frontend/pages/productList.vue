<template>
    <div>
        <div class="is-flex is-align-items-center is-justify-content-flex-end" >
            <button class="py-3 px-5 btn" @click="openModal" value="" type="submit">Добавить продукт</button>
            <AddProductForm v-if="modalVisible" @closeModal="closeModal" :modalVisible="modalVisible" />
        </div>
        <productForm @productData="handleData" />
        <ProductTable :product-list="productList" />
    </div>
</template>

<script>
import AddProductForm from '~/components/AddProductForm.vue';
import ProductTable from '~/components/ProductTable.vue';
import productForm from '~/components/ProductForm.vue';

export default {
    components: {
        productForm,
        ProductTable,
        AddProductForm
    },
    data() {
        return {
            productList: [],
            columns: [
                {
                    field: "ProductID",
                    label: "ID продукта"
                },
                {
                    field: "Name",
                    label: "Название продукта"
                },
                {
                    field: "Descr",
                    label: "Описание продукта"
                },
            ],
            modalVisible: false
        }
    },
    methods: {
        openModal() {
            this.modalVisible = true
        },
        closeModal() {
            this.modalVisible = false
        },
        handleData(productData) {
            this.productList = productData
        },
        async LoadStockList() {
            const limit = 3
            const url = `http://127.0.0.1:9000/product_list?limit=${limit}`
            try {
                const response = await fetch(url, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                })

                const responseData = await response.json()
                this.productList = responseData.Data.product_list
                console.log(this.productList)
            } catch (error) {
                console.error(error)
            }
        }

    },
    async mounted() {
        await this.LoadStockList()
    },

}
</script>