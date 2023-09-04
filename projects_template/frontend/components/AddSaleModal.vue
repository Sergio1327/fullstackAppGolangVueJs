<template>
    <div class="is-flex is-align-items-center is-justify-content-flex-end">
        <button class="btn" type="submit" @click="openModal">Добавить продажу</button>
        <AddSaleForm v-if="modalVisible" :stockOptions="stockList" :variantOptions="variantIDs" :modalVisible="modalVisible"
            @closeModal="closeModal" />
    </div>
</template>

<script>
import AddSaleForm from './AddSaleForm.vue';
export default {
    components: {
        AddSaleForm
    },
    data() {
        return {
            modalVisible: false,
            stockList: [],
            variantIDs: []
        }
    },
    methods: {
        async openModal() {
            try {

                const response = await fetch("http://127.0.0.1:9000/stock_list", {
                    method: "GET",
                    headers: {
                        'Content-Type': 'application/json'
                    },

                })
                const responseData = await response.json()
                const data = responseData.Data.stock_list
                console.log(data)
                this.stockList = data.map(e => {
                    return {
                        Option: e.StorageID,
                        Value: e.StorageID
                    }
                })

                const response2 = await fetch("http://127.0.0.1:9000/product_list?limit=999999", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json"
                    }
                })

                const responseData2 = await response2.json()
                responseData2.Data.product_list.forEach(product => {
                    product.VariantList.forEach(variant => {
                        this.variantIDs.push({
                            Option: variant.variant_id,
                            Value: variant.variant_id
                        });
                    });
                });
                this.variantIDs.sort((a, b) => a.Option - b.Option)
                this.modalVisible = true;

            } catch (error) {
                console.error(error)
            }





        },

        closeModal() {
            this.variantIDs = []
            this.modalVisible = false;
        },
    }
}
</script>

<style>
.btn {
    padding: 14px 28px;
    background-color: rgb(255, 190, 70) !important;
    border-radius: 5px;
}
</style>