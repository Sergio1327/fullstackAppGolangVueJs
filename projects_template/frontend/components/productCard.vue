<template>
    <div class="card ">
        <div class="card-content is-flex is-align-items-center is-justify-content-space-between ">
            <div class="content column is-text-center">
                <div>ID продукта : <strong>{{ productInfo.ProductID }}</strong> </div>
                <div>Название продукта : <strong>{{ productInfo.Name }}</strong></div>
                <div>Описание продукта : <strong>{{ productInfo.Descr }}</strong> </div>
                <div class="btns mt-3 is-flex is-justify-content-flex-start ">
                    <button class="button is-link" @click="viewDetails">
                        {{ showDetails ? 'Скрыть детали' : 'Просмотр деталей' }}
                    </button>
                    <button @click="toggleModal" class="button is-link ml-4 ">Добавить цены</button>
                    <button @click="stockModalToggle" class="button is-link ml-4 ">Добавить продукт на склад</button>
                </div>
                <div class="details py-5" v-for="v in variantList" :key="v.variant_id">
                    <div>ID варианта : {{ v.variant_id }}</div>
                    <div>Объем : {{ v.weight }}</div>
                    <div>Единица измерения : {{ v.unit }}</div>
                    <div>Цена : {{ v.price }}</div>
                    <div>В каких складах есть продукт :</div>
                    <div v-for="st in v.in_storages" :key="st.StorageID">
                        <div>Название склада : {{ st.StorageName }}</div>
                    </div>
                </div>
            </div>
            <div class="image is-128x128">
                <img src="@/assets/productList.jpg" class="is-img" alt="">
            </div>
        </div>
        <priceModal :showPriceModal="showPriceModal" @closeModal="closePriceModal" />
        <AddInStockModal :showStockModal="showStockModal" @closeModal="closeStockModal" />
    </div>
</template>
    
<script>
import priceModal from './priceModal.vue';
import AddInStockModal from './AddInStockModal.vue';
export default {
    components: {
        priceModal,
        AddInStockModal
    },
    props: {
        productInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            variantList: [],
            showDetails: false,
            showPriceModal: false,
            showStockModal: false,
            variantID: 0,
            startDate: null,
            price: 0,
            resp: ""
        }
    },
    methods: {
        async viewDetails() {
            if (this.showDetails) {
                this.variantList = [];
                this.showDetails = false;
            } else {
                const url = `http://127.0.0.1:9000/product/${this.productInfo.ProductID}`;
                try {
                    const response = await fetch(url, {
                        method: "GET",
                        headers: {
                            "Content-Type": "application/json",
                        },
                    });

                    const responseData = await response.json();
                    const data = responseData.Data.product_info.VariantList;
                    this.variantList = data;
                    this.showDetails = true;
                } catch (error) {
                    console.error(error);
                }
            }
        },
        toggleDetails() {
            this.showDetails = !this.showDetails
        },
        toggleModal() {
            this.showPriceModal = !this.showPriceModal
        },
        stockModalToggle() {
            this.showStockModal = !this.showStockModal
        },
        closePriceModal() {
            this.showPriceModal = false;
            this.variantID = 0;
            this.resp = '';
        },
        closeStockModal(){
            this.showStockModal = false;
            this.variantID = 0;          
            this.resp = '';
        }
    }
};
</script>
    
<style scoped>
.image {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.details {
    border-bottom: 2px solid teal;
}

.card {
    margin-bottom: 20px;
    border-radius: 5px;
    box-shadow: 1px 1px 8px orange;
}

.card-content {
    padding: 1rem;
}

.content>div>strong {
    font-weight: 600;
    color: #000;
}

.title.is-4 {
    margin-bottom: 5px;
}

.subtitle.is-6 {
    color: #999;
}
</style>