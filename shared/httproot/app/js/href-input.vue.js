import {nextTick} from 'vue'

export default {
    template: `
        <div>
            <div>
                <div class="row p-0 m-0">
                    <div class="col p-0">
                        <page-select ref="pageSelect" v-if="viewMode == 0" v-model="pageId" @update:modelValue="pageSelectChanged" :label="label"></page-select>
                        <asset-select ref="assetSelect" v-if="viewMode == 1" v-model="href" @update:modelValue="hrefChanged" :label="label"></asset-select>
                    </div>
                    <div class="col-auto p-0">
                        <button class="btn border" @click="toggleViewMode">
                            <i v-if="viewMode == 0" class="bi bi-image"></i>
                            <i v-if="viewMode == 1" class="bi bi-folder-fill"></i>
                        </button>
                    </div>
                </div>
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
                this.href = "" + newValue;
                this.updatePageId();
                this.viewMode = this.pageId ? 0 : 1;
            },
            immediate: true
        }
    },

    data() {
        return {
            href: "",
            pageId: 0,
            viewMode: 0
        }
    },

    methods: {

        toggleViewMode() {
            this.viewMode = this.viewMode < 1 ? this.viewMode + 1 : 0;
            nextTick(() => {
                if (this.viewMode === 0) {
                    this.$refs.pageSelect.showDropdown()
                } else if (this.viewMode === 1) {
                    this.$refs.assetSelect.selectAsset()
                }
            });
        },

        hrefChanged() {
            this.updatePageId()
            this.$emit('update:modelValue', this.href);
        },

        pageSelectChanged() {
            if (this.pageId)
                this.href = "page://" + this.pageId;
            else
                this.href = "";
            this.$emit('update:modelValue', this.href);
        },

        updatePageId() {
            if (this.href && this.href.length > 7 && this.href.indexOf("page://") === 0) {
                this.pageId = parseInt(this.href.substring(7).trim());
            } else {
                this.pageId = 0;
            }
        }
    },

    emits: ['update:modelValue'],
}