<template>
    <div>
        <button class="py-3 px-5 btn" @click="openModal" value="" type="submit">Добавить продукт</button>
        <AddProductForm v-if="modalVisible" @closeModal="closeModal" :modalVisible="modalVisible" />
    </div>
</template>

<script>
import AddProductForm from './AddProductForm.vue';
export default {
    components: {
        AddProductForm
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