<script>
import formVue from "~/components/SaleForm.vue";
import saleList from "~/components/saleList.vue";
import btnModal from "~/components/btnModal.vue";

export default {
    components: {
        formVue,
        saleList,
        btnModal
    }, methods: {
        handleData(data) {
            this.saleListData = data
        }
    }, data() {
        return {
            saleListData: []
        }
    }
}
</script>

<template>
    <div>
        <h1 class="label is-text-center mb-5">Продажи</h1>
        <btnModal />
        <formVue @data="handleData" />
        <saleList :saleList="saleListData" />
    </div>
</template>