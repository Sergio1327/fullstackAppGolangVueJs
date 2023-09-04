<template>
    <div>
        <h1 class="label is-text-center mb-5">Продажи</h1>

        <div class="is-flex is-align-items-center is-justify-content-flex-end">
            <button class="btn" type="submit" @click="openModal">Добавить продажу</button>
            <AddSaleForm v-if="modalVisible" :stockOptions="stockList" :variantOptions="variantIDs"
                :modalVisible="modalVisible" @closeModal="closeModal" />

        </div>

        <formVue @data="handleData" />
        <b-table class="mt-6" :data="saleListData" :columns="columns"></b-table>
    </div>
</template>



<script>
import formVue from "~/components/SaleForm.vue";
import AddSaleForm from "~/components/AddSaleForm.vue";
export default {
    components: {
        formVue,
        AddSaleForm
    }, methods: {
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

        handleData(data) {
            this.saleListData = data
        },

        formatDate(dateTime) {
            const date = new Date(dateTime);
            const year = date.getUTCFullYear();
            const month = String(date.getUTCMonth() + 1).padStart(2, "0");
            const day = String(date.getUTCDate()).padStart(2, "0");
            const hours = String(date.getHours()).padStart(2, "0");
            const minutes = String(date.getUTCMinutes()).padStart(2, "0");
            const seconds = String(date.getUTCSeconds()).padStart(2, "0");
            return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}+05:00`;
        },

        async LoadSales() {
            try {
                const formattedStartDate = this.formatDate(this.startDate);
                const formattedEndDate = this.formatDate(this.endDate);

                const requestData = {
                    start_date: formattedStartDate,
                    end_date: formattedEndDate,
                    limit: parseInt(this.limit),
                    product_name: this.productName,
                    storage_id: parseInt(this.storageId)
                };
                const response = await fetch("http://127.0.0.1:9000/sales", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData)
                })

                const responseData = await response.json()
                this.saleListData = responseData.Data.sale_list

            } catch (error) {
                console.error(error)
            }
        }

    }, data() {
        return {
            modalVisible: false,
            stockList: [],
            variantIDs: [],

            startDate: new Date("2020-08-20T00:00:00"),
            endDate: new Date("2029-08-20T00:00:00"),
            limit: 1,
            productName: "",
            storageId: 1,

            saleListData: [],

            columns: [
                {
                    field: "SaleID",
                    label: "ID Продажи"
                },
                {
                    field: "ProductName",
                    label: "Название продукта"
                },
                {
                    field: "variant_id",
                    label: "ID варианта"
                },
                {
                    field: "storage_id",
                    label: "ID склада"
                },
                {
                    field: "SoldAt",
                    label: "дата продажи"
                },
                {
                    field: "quantity",
                    label: "Колличество продукта"
                },
                {
                    field: "TotalPrice",
                    label: "Общая цена продажи"
                }

            ]
        }
    },
    async mounted() {
        await this.LoadSales()
    }
}
</script>

