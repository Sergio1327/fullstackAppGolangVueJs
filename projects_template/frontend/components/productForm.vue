<template>
    <div>
        <div class="field">
            <label class="label">Название продукта</label>
            <div class="control">
                <input class="input" type="text" v-model="productName">
            </div>
        </div>
        <div class="field">
            <label class="label">Тег продукта</label>
            <div class="control">
                <input class="input" type="text" v-model="tag">
            </div>
        </div>
        <div class="field">
            <label class="label">Лимит вывода</label>
            <div class="control">
                <div class="select">
                    <select v-model="limit">
                        <option value="1">1</option>
                        <option value="5">5</option>
                        <option value="10">10</option>
                        <option value="50">50</option>
                    </select>
                </div>
            </div>
        </div>
        <button class="button is-link" @click="GetProductList">Найти список продуктов</button>
    </div>
</template>
<script>
export default {
    data() {
        return {
            productName: "",
            tag: "",
            limit: 1
        }
    },
    methods: {
        async GetProductList() {
          const  limitProductList = parseInt(this.limit)
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
                this.$emit("productData", productData)

            } catch (error) {
                console.error(error)
            }


        }
    }
}
</script>