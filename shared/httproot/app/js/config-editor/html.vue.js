export default {
    template: `
        <div>
            <neiki-editor v-model="editorContent"></neiki-editor>
        </div>
    `,

    data() {
        return {
            editorContent: ""
        }
    },

    methods: {
        setValue(value) {
            this.editorContent = value ? value : "&nbsp;"; //TODO find why editor not showing if no value
        },

        getValue() {
            return this.editorContent;
        },
    },
}