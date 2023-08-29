<script>
import formVue from "~/components/SaleForm.vue";
import btnModal from "~/components/BtnModal.vue";

export default {
    components: {
        formVue,
        btnModal
    }, methods: {
        handleData(data) {
            this.saleListData = data
        }
    }, data() {
        return {
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
    }
}
</script>

<template>
    <div>
        <h1 class="label is-text-center mb-5">Продажи</h1>
        <btnModal />
        <formVue @data="handleData" />
        <b-table class="mt-6" :data="saleListData" :columns="columns"></b-table>
    </div>
</template>