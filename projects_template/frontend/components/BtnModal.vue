<template>
    <div class="is-flex is-align-items-center is-justify-content-flex-end">
        <button class="btn" type="submit" @click="openModal">Добавить продажу</button>
        <addSaleForm v-if="modalVisible" :modalVisible="modalVisible" @closeModal="closeModal" />
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
        }
    },
    methods: {
        openModal() {
            this.modalVisible = true;
        },
        submitModalData() {
            this.modalVisible = false;
        },
        closeModal() {
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