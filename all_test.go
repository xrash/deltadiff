package deltadiff

import (
	"bytes"
	"fmt"
	"github.com/franela/goblin"
	"github.com/xrash/deltadiff/testdata"
	"testing"
	"time"
)

func TestAll(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("deltadiff", func() {

		g.It("should work in one binary case", func() {
			g.Timeout(time.Second * 60)

			blocksizes := []int{
				256,
				512,
				1024,
				2048,
			}

			base, ok1 := testdata.FS.Get("/maamoul-mod.jpg")
			target, ok2 := testdata.FS.Get("/maamoul.jpg")
			g.Assert(ok1).Equal(true)
			g.Assert(ok2).Equal(true)

			dc := &DeltaConfig{
				Debug: false,
			}

			for _, blocksize := range blocksizes {
				sc := &SignatureConfig{
					Hasher:    "polyroll",
					BlockSize: blocksize,
					BaseSize:  len(base),
				}

				baseBuffer := bytes.NewBuffer(base)
				signatureBuffer := bytes.NewBuffer(nil)
				errSignature := Signature(baseBuffer, signatureBuffer, sc)
				g.Assert(errSignature).Equal(nil)

				targetBuffer := bytes.NewBuffer(target)
				deltaBuffer := bytes.NewBuffer(nil)
				errDelta := Delta(signatureBuffer, targetBuffer, deltaBuffer, dc)
				g.Assert(errDelta).Equal(nil)

				outBuffer := bytes.NewBuffer(nil)
				baseBuffer = bytes.NewBuffer(base)
				errPatch := Patch(baseBuffer, deltaBuffer, outBuffer)
				g.Assert(errPatch).Equal(nil)

				equal := outBuffer.String() == string(target)
				if !equal {
					fmt.Println("blocksize", blocksize)
				}
				g.Assert(equal).Equal(true)
			}
		})

		g.It("should work in several different cases", func() {

			run := func(base, target string, sc *SignatureConfig) {
				dc := &DeltaConfig{
					Debug: false,
				}

				baseBuffer := bytes.NewBufferString(base)
				signatureBuffer := bytes.NewBuffer(nil)
				errSignature := Signature(baseBuffer, signatureBuffer, sc)
				g.Assert(errSignature).Equal(nil)

				targetBuffer := bytes.NewBufferString(target)
				deltaBuffer := bytes.NewBuffer(nil)
				errDelta := Delta(signatureBuffer, targetBuffer, deltaBuffer, dc)
				g.Assert(errDelta).Equal(nil)

				outBuffer := bytes.NewBuffer(nil)
				baseBuffer = bytes.NewBufferString(base)
				errPatch := Patch(baseBuffer, deltaBuffer, outBuffer)
				g.Assert(errPatch).Equal(nil)

				g.Assert(outBuffer.String()).Equal(target)
			}

			testcases := [][]string{
				[]string{
					"aaaaaabbbbbbccccddeeedeeeeeeea",
					"aaaaaabbbbbbccccccddddddeeeeee",
				},
				[]string{
					"",
					"aaaaaabbbbbbccccccddddddeeeeee",
				},
				[]string{
					"aaaaaabbbbbbccccddeeedeeeeeeea",
					"",
				},
				[]string{
					"",
					"",
				},
				[]string{
					"dasasdjaidhasdojadsiaosd asdsad saasd$$$$ DASDAASASDADS   ",
					"dasasdjaidhasdojadsiaosd $$$$ DASDAASASDADS  WWWWWWWWWWWWW",
				},
				[]string{
					"",
					"a",
				},
				[]string{
					"a",
					"",
				},
				[]string{
					"SSenhor. ??? Eu nam escrevo a vos alteza per minha m??o, porque, quando esta fa??o, tenho muito grande salu??o, que he sinal de morrer: eu, senhor, deixo qu?? ese filho per minha mem??ria, a que deixo toda minha fazemda, que he assaz de pouca, mas deixo lhe a obriga??am de todos meus serui??os, que he mui grande: as cousas da india ellas falar??m por mim e por elle: deixo a india com as principaes cabe??as tomadas em voso poder, sem nela ficar outra penden??a senam cerrar se e mui bem a porta do estreito;a isto he o que me vosa alteza encomendou: eu, senhor, vos dey sempre por comselho, pera segurar de l?? india, irdes vos tirando de despesas: pe??oa vos alteza por mercee que se lembre de tudo isto, e que me fa??a meu filho grande, e lhe d?? toda satisfa??am de meu serui??o: todas minhas confian??as pus nas m??s de vos alteza e da senhora Rainha, a elles m emcomemwdo, que fa??am mwinhas cousas grandes, pois acabo em cousas de voso serui??o, e por elles vollo tenho merecido; e as minhas ten??as, as quaes comprey pela maior parte, como vossa alteza sabe, beijar lh ey as m??os pollas em meu filho: escrita no mar a 6 dias de dezembro de 1515. Afomso dalboquerqueL",
					"Senhor. ??? Eu nam escrevo a vos alteza per minha m??o, porque, quando esta fa??o, tenho muito grande salu??o, que he sinal de morrer: eu, senhor, deixo qu?? ese filho per minha mem??ria, a que deixo toda minha fazemda, que he assaz de pouca, mas deixo lhe a obriga??am de todos meus serui??os, que he mui grande: as cousas da india ellas falar??m por mim e por elle: deixo a india com as principaes cabe??as tomadas em voso poder, sem nela ficar outra penden??a senam cerrar se e mui bem a porta do estreito; isto he o que me vosa alteza encomendou: eu, senhor, vos dey sempre por comselho, pera segurar de l?? india, irdes vos tirando de despesas: pe??oa vos alteza por mercee que se lembre de tudo isto, e que me fa??a meu filho grande, e lhe d?? toda satisfa??am de meu serui??o: todas minhas confian??as pus nas m??s de vos alteza e da senhora Rainha, a elles m emcomemdo, que fa??am minhas cousas grandes, pois acabo em cousas de voso serui??o, e por elles vollo tenho merecido; e as minhas ten??as, as quaes comprey pela maior parte, como vossa alteza sabe, beijar lh ey as m??os pollas em meu filho: escrita no mar a 6 dias de dezembro de 1515. Afomso dalboquerque",
				},
			}

			hashers := []string{
				"polyroll",
				"md5",
				"crc32",
			}

			for _, testcase := range testcases {
				for i := 1; i <= 10; i++ {
					for _, hasher := range hashers {

						base := testcase[0]
						target := testcase[1]

						signatureConfig := &SignatureConfig{
							Hasher:    hasher,
							BlockSize: i * i,
							BaseSize:  len(base),
						}

						run(base, target, signatureConfig)
					}
				}
			}
		})

	})
}
