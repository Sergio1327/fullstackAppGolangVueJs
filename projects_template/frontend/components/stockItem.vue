<template>
    <li>
        <div class="card ">
            <div class="card-content is-flex is-align-items-center is-justify-content-space-between ">
                <div class="content column is-text-center">
                    <div>ID склада : <strong> {{ stock.StorageID }}</strong></div>
                    <div>Название склада : <strong> {{ stock.StorageName }}</strong></div>
                </div>
                <div class="image is-128x128">
                    <img src="@/assets/stockList.jpg" class="is-img" alt="">
                </div>
            </div>
        </div>
    </li>
</template>


<script>
export default {
    props: {
        stock: {
            type: Object,
            required: true
        }
    }
}
</script>
<style scoped >
content div strong {
    color: #000;
    font-weight: 600;
}
</style>