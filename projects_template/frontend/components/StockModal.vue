<template>
    <div class="is-flex is-align-items-center is-justify-content-flex-end">
        <button class="btn mb-5" @click="openModal" value="" type="submit">Добавить склад</button>
        <AddStockForm v-if="modalVisible" @dataSubmitted="submitModalData" @closeModal="closeModal" />
    </div>
</template>

<script>
import AddStockForm from './AddStockForm.vue';
export default {

    components: {
        AddStockForm
    },
    data() {
        return {
            modalVisible: false,
        }
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

}
</script>

<style>
.btn {
    padding: 14px 28px;
    background-color: rgb(255, 190, 70) !important;
    border-radius: 5px;
}
</style>