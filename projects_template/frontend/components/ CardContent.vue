<template>
    <div class="card mb-5  py-3 ">
        <nuxt-link :to="card.link">
            <div class="card-content is-flex is-align-items-center">
                <div class="media m-0">
                    <div class="media-left">
                        <figure class="image is-96x96 is-flex is-align-items-center">
                            <img :src="card.imgLink" alt="Placeholder image">
                        </figure>
                    </div>
                </div>
                <div class="media-content">
                    <p class="title is-4 m-0">{{ card.title }}</p>
                    <div class="content has-text-warning">
                        {{ card.descr }}
                    </div>
                </div>

            </div>
        </nuxt-link>
    </div>
</template>

<script>

export default {
    props: {
        card: {
            type: Object,
            required: true
        }
    }
}
</script>

<style scoped   >
.card {
    transition: 200ms linear !important;
}

.card:hover {
    transform: scale(1, 1.1);
    box-shadow: 1px 1px 10px orange;
}
.content{
    color: orange !important;
}
</style>