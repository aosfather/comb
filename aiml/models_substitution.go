package aiml

//
type FindReplaceOperation interface {
	Matches(index int, input []string) bool
	Replacement(index int, output []string) int
}

type Substitution struct {
	operation FindReplaceOperation
	find      string
	replace   string
	tokenizer *Tokenizer
}

func (this *Substitution) Init(find, replace string, tokenizer *Tokenizer) {
	this.find = find
	this.replace = replace
	this.tokenizer = tokenizer
	this.afterSetProperty()

}
func (this *Substitution) afterSetProperty() {
	// if (find == null || tokenizer == null || replace == null)
	//      return;

	//    List<String> tokens = tokenizer.tokenize(find);
	//    if (tokens.size() > 1)
	//      operation = new FindReplaceFragment(tokens);
	//    else if (find.charAt(0) != ' ')
	//      operation = new FindReplaceSuffix();
	//    else if (find.charAt(find.length() - 1) != ' ')
	//      operation = new FindReplacePrefix();
	//    else
	//      operation = new FindReplaceWord();

	//    find = find.toUpperCase().trim();
}

func (this *Substitution) Substitute(input []string) {

}

func (this *Substitution) SubstituteOffset(offset int, input []string) int {

}

type FindReplaceFragment struct {
	owner       Substitution
	replacement []string
	fragment    []string
}

type FindReplacePrefix struct {
	owner   Substitution
	token   string
	upToken string
}

type FindReplaceSuffix struct {
	owner   Substitution
	token   string
	upToken string
}

type FindReplaceWord struct {
	owner       Substitution
	replacement []string
}
