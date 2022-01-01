package hasher

import (
	"github.com/franela/goblin"
	"github.com/xrash/deltadiff/testdata"
	"testing"
	"time"
)

var teststrings = []string{

	"",

	"a",

	"aa",

	"aaa",

	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",

	`
e porque Vossa Alteza me disse que se queria
nisto fiar de mim, portanto preparei fazer um livro de
cosmografia e marinharia, cujo prólogo é este que aqui
é escrito, o qual livro será partido em cinco livros, e no
primeiro se dirá do que descobrio o virtuoso Infante
Dom Anrique, e no segundo do que mandou descobrir o
excelente rei Dom Afonso, e no terceiro do que isso
mesmo fez descobrir o sereníssimo rei Dom João, que
faz fim no ilhéo da Cruz, como já disse; o quarto e o
quinto, em que pendem vossos gloriosos feitos, que são
mais em cantidade e maiores em calidade que os de tôdolos outros príncipes.
`,

	`
Varre, varre, varre, varre vassourinha!
Varre, varre a bandalheira!
Que o povo já 'tá cansado
De sofrer dessa maneira
Jânio Quadros é a esperança desse povo abandonado!
Jânio Quadros é a certeza de um Brasil, moralizado!
Alerta, meu irmão!
Vassoura, conterrâneo!
Vamos vencer com Jânio!
`,

	`
Varre, varre, varre, varre vassourinha!
Varre, varre a bandalheira!
Que o povo já 'tá cansado
De sofrer dessa maneira
Jânio Quadros é a esperança desse povo abandonado!
Jânio Quadros é a certeza de um Brasil, moralizado!
Alerta, meu irmão!
Vassoura, conterrâneo!
Vamos vencer com Jânio!
`,
}

var basemodlist = [][]int{
	[]int{7, 113},
	[]int{257, 4909},
	[]int{23, 709},
}

func TestPolyrollHasher(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("polyroll hasher", func() {

		g.It("some hardcoded tests", func() {
			h := &PolyrollHasher{
				Base: 7,
				Mod:  113,
			}

			inputs := []string{
				"abcd",
				"bcde",
			}

			expecteds := [][]byte{
				[]byte{0, 0, 0, 107},
				[]byte{0, 0, 0, 55},
			}

			for i := 0; i < len(inputs); i++ {
				input := inputs[i]
				expected := expecteds[i]
				hash, err := h.Hash([]byte(input))
				g.Assert(err).Equal(nil)
				g.Assert(string(hash)).Equal(string(expected))
			}
		})

		g.It("hashing individually and rollingly should be equal (binary)", func() {
			g.Timeout(time.Second * 60)

			input1, ok1 := testdata.FS.Get("/maamoul-mod.jpg")
			input2, ok2 := testdata.FS.Get("/maamoul.jpg")
			g.Assert(ok1).Equal(true)
			g.Assert(ok2).Equal(true)

			inputs := [][]byte{
				input1,
				input2,
			}

			blockSizes := []int{
				256,
				512,
				1024,
			}

			for _, input := range inputs {
				for _, blockSize := range blockSizes {
					for i := 0; i < len(input)-blockSize; i++ {

						h := &PolyrollHasher{
							Base: POLYROLL_BASE,
							Mod:  POLYROLL_MOD,
						}

						segment := input[i : i+blockSize]

						hash, err := h.Hash(segment)
						g.Assert(err).Equal(nil)
						singleHash, err := h.SingleHash(segment)
						g.Assert(err).Equal(nil)

						g.Assert(string(hash)).Equal(string(singleHash))
					}
				}
			}
		})

		g.It("hashing individually and rollingly should be equal (strings)", func() {
			for _, basemod := range basemodlist {
				for blockSize := 1; blockSize < 23; blockSize++ {
					for i := 0; i < len(teststrings); i++ {
						input := []byte(teststrings[i])

						h := &PolyrollHasher{
							Base: basemod[0],
							Mod:  basemod[1],
						}

						for i := 0; i < len(input)-blockSize; i++ {
							segment := input[i : i+blockSize]

							hash, err := h.Hash(segment)
							g.Assert(err).Equal(nil)
							singleHash, err := h.SingleHash(segment)
							g.Assert(err).Equal(nil)

							g.Assert(string(hash)).Equal(string(singleHash))
						}
					}
				}
			}

		})

	})
}
