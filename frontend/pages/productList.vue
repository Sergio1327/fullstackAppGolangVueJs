<template>
    <section>
        <div class="is-flex is-align-items-center is-justify-content-flex-end">
            <button class="button is-warning py-3 px-5" @click="openModal" value="" type="submit">Добавить продукт</button>
            <AddProductForm v-if="modalVisible" @closeModal="closeModal" :modalVisible="modalVisible" />
        </div>
        <div>
            <div class="field">
                <label class="label">Название продукта</label>
                <div class="control">
                    <input class="input" placeholder="Введите название продукта" type="text" v-model="productName">
                </div>
            </div>
            <div class="field">
                <label class="label">Тег продукта</label>
                <div class="control">
                    <input class="input" placeholder="Введите тег продукта" type="text" v-model="tag">
                </div>
            </div>
            <div class="field">
                <label class="label">Лимит вывода</label>
                <div class="control">
                    <div class="select">
                        <b-select placeholder="введите лимит вывода"  v-model="limit">
                            <option value="1">1</option>
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="50">50</option>
                        </b-select>
                    </div>
                </div>
            </div>
            <button class="button is-link" @click="GetProductList">Найти список продуктов</button>
        </div>
        <ProductTable :product-list="productList" />
    </section>
</template>

<script>
import AddProductForm from '~/components/AddProductForm.vue';
import ProductTable from '~/components/ProductTable.vue';

export default {
    components: {
        ProductTable,
        AddProductForm
    },

    data() {
        return {
            productName: "",
            tag: "",
            limit: null,
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
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        },

        async GetProductList() {
            const limitProductList = parseInt(this.limit)
            const url = `http://127.0.0.1:9000/product_list?name=${this.productName}&tag=${this.tag}&limit=${limitProductList}`

            try {
                const response = await fetch(url, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                })

                const responseData = await response.json()
                const productData = responseData.Data.product_list
                this.productList = productData
                this.productName = ""
                this.tag = ""

            } catch (error) {
                this.$buefy.snackbar.open({
                    message: `${error}`,
                    type: "is-danger"
                })
                console.error(error)
            }
        }
    },

    async mounted() {
        await this.LoadStockList()
    },

}
</script>

<style scoped >
.input {
    width: 30% !important;
}

@media (max-width:1000px) {
    .input {
        width: 100% !important;
    }
}
</style>