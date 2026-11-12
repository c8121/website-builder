export default {
    template: `
        <div>
            <template v-if="editorContent">
                <textarea style="height: 50vh" v-model="editorContent.text" class="form-control"></textarea>
            </template>
        </div>
    `,

    data() {
        return {
            editorContent: null
        }
    },


    methods: {
        setContent(content) {
            if (content) {
                try {
                    this.editorContent = this.initContent(JSON.parse(content));
                } catch (e) {
                    console.log("Failed to parse content");
                    this.editorContent = this.initContent({text: content});
                }
            } else {
                this.editorContent = this.initContent({});
            }
        },

        initContent(content) {
            if (!content.text) {
                content.text = "";
            }

            return content;
        },

        getContent() {
            return JSON.stringify(this.editorContent);
        },
    },

    emits: ['componentEvent'],
}