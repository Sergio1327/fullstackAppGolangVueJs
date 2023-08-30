<template>
    <section>
        <b-modal v-model="isActive" has-modal-card trap-focus :destroy-on-hide="false" aria-role="dialog"
            aria-label="Добавление склада" close-button-aria-label="Закрыть" aria-modal>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">Варианты продукта</p>
                    <button type="button" class="delete" @click="closeModal"></button>
                </header>
                <section class="modal-card-body">
                    <b-field class="is-flex is-flex-direction-column">
                        <div class="variant is-flex is-flex-direction-column p-2" v-for="v in variantList"
                            :key="v.variant_id">
                            <div>ID варианта : {{ v.variant_id }}</div>
                            <div>Объем : {{ v.weight }}{{ v.unit }} </div>
                            <div>Цена : {{ v.price }}$</div>
                            <div class="storages">В каких складах имеется :
                                <div v-for="s in v.in_storages" :key="s.storageID">
                                    {{ s.StorageName }}
                                </div>
                            </div>
                        </div>
                    </b-field>
                </section>
                <footer class="modal-card-foot">
                    <b-button label="Закрыть" @click="closeModal" />
                </footer>

            </div>
        </b-modal>
    </section>
</template>

<script>
export default {
    data() {
        return {
            isActive: this.modalVisible,
        }
    },
    props: {
        modalVisible: {
            type: Boolean,
            required: true
        },
        variantList: {
            type: Array,
            required: true
        }
    },
    methods: {
        closeModal() {
            this.$emit("closeModal")
        }
    }
}
</script>


<style>
.variant {
    border-bottom: 2px solid teal;
}

.field .has-addons {
    display: flex !important;
    flex-direction: column !important;
}

.storages {
    width: max-content !important;
}
</style>    