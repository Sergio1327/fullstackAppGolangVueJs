<template>
    <div>
        <button class="py-3 px-5 btn" @click="openModal" value="" type="submit">Добавить продукт</button>
        <productModal v-if="modalVisible" @dataSubmitted="submitModalData" @closeModal="closeModal" />
    </div>
</template>

<script>
import productModal from './productModal.vue';
export default {
    components: {
        productModal
    }, methods: {
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
            modalVisible: false,
        }
    }
}
</script>

<style scoped></style>