<script>
import formVue from "~/components/SaleForm.vue";
import AddSaleModal from "~/components/AddSaleModal.vue";

export default {
    components: {
        formVue,
        AddSaleModal
    }, methods: {
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

<template>
    <div>
        <h1 class="label is-text-center mb-5">Продажи</h1>
        <AddSaleModal />
        <formVue @data="handleData" />
        <b-table class="mt-6" :data="saleListData" :columns="columns"></b-table>
    </div>
</template>