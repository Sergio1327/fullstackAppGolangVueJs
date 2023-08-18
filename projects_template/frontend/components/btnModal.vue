<template>
    <div class="is-flex is-align-items-center is-justify-content-flex-end">
        <button class="btn" @click="openModal" value="" type="submit">Добавить продажу</button>
        <SaleModal v-if="modalVisible" @dataSubmitted="submitModalData" @closeModal="closeModal" />
    </div>
</template>

<script>
import SaleModal from './SaleModal.vue';
export default {
    components: {
        SaleModal
    },
    methods: {
        openModal() {
            this.modalVisible = true;
        },
        submitModalData() {
            // Здесь можно обновить список данных
            this.modalVisible = false;
        },
        closeModal() {
            this.modalVisible = false;
        },
    },
    data() {
        return {
            saleListData: [],
            modalVisible: false,
        }
    }
}
</script>

<style>
.btn{
    padding: 14px 28px;
    background-color: rgb(255, 190, 70) !important;
    border-radius: 5px;
}
</style>