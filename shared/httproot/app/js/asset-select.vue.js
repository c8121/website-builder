import AssetStorageMain from './asset-storage-main.vue.js'
import Dialog from './dialog.js'

export default {
    template: `
        <div>
            <div v-if="url" class="form-control p-2 mb-1 overflow-scroll">
                <img :src="url">
            </div>
            <div class="input-group">
                <label v-if="label" class="input-group-text">{{label}}</label>
                <input class="form-control" v-model="url">
                <button class="btn border" @click="selectAsset"><i class="bi bi-folder2-open"></i></button>
            </div>
        </div>
    `,

    props: {
        modelValue: {
            type: String,
            default: null
        },
        label: {
            type: String,
            default: null
        }
    },

    watch: {
        modelValue: {
            handler(newValue) {
                this.url = newValue;
            },
            immediate: true
        }
    },

    data() {
        return {
            url: ""
        }
    },

    methods: {

        selectAsset() {
            const self = this;

            const response = fetch(new URL('./asset-storage-main.tpl.html', import.meta.url))
                    .then(response => response.text())
                    .then(html => {
                        AssetStorageMain.template = html;


                        const dialog = Dialog.create();
                        dialog.setTitle("Test dialog");
                        dialog.setContentComponent(AssetStorageMain, {
                            downloadCallback: (url) => {
                                console.log("Asset selected", url);
                                this.$emit('update:modelValue', url);
                                dialog.close();
                            }
                        });        
                    });
        }
    },

    emits: ['update:modelValue'],
}