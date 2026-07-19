package i18n

// Lang is a language code.
type Lang string

const (
	EN    Lang = "en"
	ZH_CN Lang = "zh-CN"
	JA    Lang = "ja"
	KO    Lang = "ko"
	ES    Lang = "es"
	PT_BR Lang = "pt-BR"
	DE    Lang = "de"
	FR    Lang = "fr"
)

// AllLangs returns all supported languages in display order.
func AllLangs() []Lang {
	return []Lang{EN, ZH_CN, JA, KO, ES, PT_BR, DE, FR}
}

// LangName returns the native name of a language.
func LangName(l Lang) string {
	names := map[Lang]string{
		EN:    "English",
		ZH_CN: "简体中文",
		JA:    "日本語",
		KO:    "한국어",
		ES:    "Español",
		PT_BR: "Português (BR)",
		DE:    "Deutsch",
		FR:    "Français",
	}
	if n, ok := names[l]; ok {
		return n
	}
	return string(l)
}

// DefaultStyleForLang returns the recommended style preset for a language.
func DefaultStyleForLang(l Lang) string {
	switch l {
	case ZH_CN:
		return "cn-official"
	case JA:
		return "jp-formal"
	case KO:
		return "kr-standard"
	default:
		return "us-business"
	}
}

// ---------------------------------------------------------------------------
// Translation maps — keyed by Lang
// ---------------------------------------------------------------------------

// T returns the translation for a key in the given language.
func T(l Lang, key string) string {
	if m, ok := dict[l]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	// Fallback to English
	if m, ok := dict[EN]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}

// dict holds all translation maps.
var dict = map[Lang]map[string]string{
	EN:    en,
	ZH_CN: zhCN,
	JA:    ja,
	KO:    ko,
	ES:    es,
	PT_BR: ptBR,
	DE:    de,
	FR:    fr,
}

// PresetDescription returns a localized description for a preset name.
func PresetDescription(l Lang, preset string) string {
	key := "preset_" + preset
	if v := T(l, key); v != key {
		return v
	}
	return preset
}

// ---------------------------------------------------------------------------
// English
// ---------------------------------------------------------------------------

var en = map[string]string{
	// App
	"app_title":       "md2docx — Markdown to DOCX",
	"app_short":       "Convert Markdown to professional DOCX documents",

	// Language selection
	"lang_select_title":   "Select Language / 选择语言",
	"lang_select_help":    "↑↓ navigate   enter select   q quit",
	"lang_saved":          "Language saved: %s",
	"lang_saved_hint":     "Change with: md2docx --lang <code>\n",

	// TUI steps
	"tui_input_title":     "md2docx — Select Markdown Input",
	"tui_output_title":    "md2docx — Choose Output",
	"tui_style_title":     "md2docx — Choose Style",
	"tui_style_custom":    "Custom JSON file...",
	"tui_style_default":   "Use default",
	"tui_style_file_title":"md2docx — Select Style Template (JSON)",
	"tui_confirm_title":   "md2docx — Confirm Conversion",
	"tui_confirm_input":   "Input",
	"tui_confirm_output":  "Output",
	"tui_confirm_style":   "Style",
	"tui_confirm_mermaid": "Render mermaid diagrams",
	"tui_confirm_convert": "Convert",
	"tui_confirm_back":    "Back",
	"tui_converting":      "Converting...",
	"tui_done":            "Done",
	"tui_done_msg":        "Done: %s\n\nPress any key to exit.",
	"tui_error_msg":       "Error: %v\n\nPress esc to exit.",
	"tui_filename_label":  "Filename",
	"tui_nav_help":        "↑↓ navigate • enter select • esc back",
	"tui_nav_quit":        "↑↓ navigate • enter select • esc quit",
	"tui_nav_tab":         "↑↓ navigate • enter confirm • tab toggle input • esc back",

	// CLI output
	"cli_ok":             "OK",
	"cli_error":          "ERROR",
	"cli_bytes":          "bytes",
	"cli_template_created": "Style template created: %s",
	"cli_no_agents":        "No agents detected.",

	// Preset descriptions
	"preset_us-business": "US Business – Cambria/Calibri, blue accent",
	"preset_us-modern":   "US Modern – Segoe UI, minimal dark tones",
	"preset_cn-official": "CN Official – SimHei/SimSun (公文风格), red accent",
	"preset_cn-modern":   "CN Modern – Noto Sans SC, modern Chinese",
	"preset_jp-formal":   "JP Formal – Yu Mincho/Yu Gothic, business",
	"preset_eu-clean":    "EU Clean – Helvetica/Arial, minimalist",
	"preset_kr-standard": "KR Standard – Malgun Gothic/Nanum, clean",
	"preset_academic":    "Academic – Times New Roman, scholarly",
	"preset_default":     "Default – Aptos Display/Cascadia Mono",

	// Error messages
	"err_markdown_not_found":   "Markdown input not found",
	"err_output_ext":           "Output path must end in .docx",
	"err_style_not_found":      "Style template not found",
	"err_style_invalid_json":   "Style template is not valid JSON",
	"err_style_missing_field":  "Style template is missing field",
	"err_style_bad_color":      "Color must be a #RRGGBB value",
	"err_style_bad_size":       "Size must be positive",
	"err_reading_file":         "Error reading file",
	"err_writing_file":         "Error writing file",
	"err_converting":           "Conversion failed",
	"err_tui":                  "TUI error",
}

// ---------------------------------------------------------------------------
// 简体中文 (Simplified Chinese)
// ---------------------------------------------------------------------------

var zhCN = map[string]string{
	"app_title":       "md2docx — Markdown 转 DOCX",
	"app_short":       "将 Markdown 转换为专业 DOCX 文档",

	"lang_select_title":   "选择语言 / Select Language",
	"lang_select_help":    "↑↓ 导航   Enter 选择   Q 退出",
	"lang_saved":          "语言已保存：%s",
	"lang_saved_hint":     "可通过 md2docx --lang <code> 修改\n",

	"tui_input_title":     "md2docx — 选择 Markdown 输入",
	"tui_output_title":    "md2docx — 选择输出位置",
	"tui_style_title":     "md2docx — 选择样式",
	"tui_style_custom":    "自定义 JSON 文件...",
	"tui_style_default":   "使用默认样式",
	"tui_style_file_title":"md2docx — 选择样式模板 (JSON)",
	"tui_confirm_title":   "md2docx — 确认转换",
	"tui_confirm_input":   "输入",
	"tui_confirm_output":  "输出",
	"tui_confirm_style":   "样式",
	"tui_confirm_mermaid": "渲染 Mermaid 图表",
	"tui_confirm_convert": "转换",
	"tui_confirm_back":    "返回",
	"tui_converting":      "转换中...",
	"tui_done":            "完成",
	"tui_done_msg":        "完成：%s\n\n按任意键退出。",
	"tui_error_msg":       "错误：%v\n\n按 Esc 退出。",
	"tui_filename_label":  "文件名",
	"tui_nav_help":        "↑↓ 导航 • Enter 选择 • Esc 返回",
	"tui_nav_quit":        "↑↓ 导航 • Enter 选择 • Esc 退出",
	"tui_nav_tab":         "↑↓ 导航 • Enter 确认 • Tab 切换输入 • Esc 返回",

	"cli_ok":             "成功",
	"cli_error":          "错误",
	"cli_bytes":          "字节",
	"cli_template_created": "样式模板已创建：%s",
	"cli_no_agents":        "未检测到 Agent。",

	"preset_us-business": "美国商务 – Cambria/Calibri，蓝色强调",
	"preset_us-modern":   "美国现代 – Segoe UI，简约深色",
	"preset_cn-official": "中国公文 – 黑体/宋体（公文风格），红色强调",
	"preset_cn-modern":   "中国现代 – Noto Sans SC，现代中文",
	"preset_jp-formal":   "日本正式 – Yu Mincho/Yu Gothic，商务",
	"preset_eu-clean":    "欧洲简约 – Helvetica/Arial，极简",
	"preset_kr-standard": "韩国标准 – Malgun Gothic/Nanum，清爽",
	"preset_academic":    "学术 – Times New Roman，学术风格",
	"preset_default":     "默认 – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "找不到 Markdown 输入文件",
	"err_output_ext":           "输出路径必须以 .docx 结尾",
	"err_style_not_found":      "找不到样式模板",
	"err_style_invalid_json":   "样式模板不是有效的 JSON",
	"err_style_missing_field":  "样式模板缺少字段",
	"err_style_bad_color":      "颜色必须是 #RRGGBB 格式",
	"err_style_bad_size":       "字号必须为正数",
	"err_reading_file":         "读取文件出错",
	"err_writing_file":         "写入文件出错",
	"err_converting":           "转换失败",
	"err_tui":                  "TUI 错误",
}

// ---------------------------------------------------------------------------
// 日本語 (Japanese)
// ---------------------------------------------------------------------------

var ja = map[string]string{
	"app_title":       "md2docx — Markdown → DOCX",
	"app_short":       "MarkdownをプロフェッショナルなDOCX文書に変換",

	"lang_select_title":   "言語選択 / Select Language",
	"lang_select_help":    "↑↓ 移動   Enter 選択   Q 終了",
	"lang_saved":          "言語を保存しました：%s",
	"lang_saved_hint":     "変更: md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Markdown入力を選択",
	"tui_output_title":    "md2docx — 出力先を選択",
	"tui_style_title":     "md2docx — スタイルを選択",
	"tui_style_custom":    "カスタムJSONファイル...",
	"tui_style_default":   "デフォルトを使用",
	"tui_style_file_title":"md2docx — スタイルテンプレートを選択 (JSON)",
	"tui_confirm_title":   "md2docx — 変換を確認",
	"tui_confirm_input":   "入力",
	"tui_confirm_output":  "出力",
	"tui_confirm_style":   "スタイル",
	"tui_confirm_mermaid": "Mermaid図を描画",
	"tui_confirm_convert": "変換",
	"tui_confirm_back":    "戻る",
	"tui_converting":      "変換中...",
	"tui_done":            "完了",
	"tui_done_msg":        "完了：%s\n\n任意のキーで終了。",
	"tui_error_msg":       "エラー：%v\n\nEscで終了。",
	"tui_filename_label":  "ファイル名",
	"tui_nav_help":        "↑↓ 移動 • Enter 選択 • Esc 戻る",
	"tui_nav_quit":        "↑↓ 移動 • Enter 選択 • Esc 終了",
	"tui_nav_tab":         "↑↓ 移動 • Enter 確認 • Tab 切替 • Esc 戻る",

	"cli_ok":             "成功",
	"cli_error":          "エラー",
	"cli_bytes":          "バイト",
	"cli_template_created": "スタイルテンプレート作成：%s",
	"cli_no_agents":        "エージェントが検出されませんでした。",

	"preset_us-business": "米国ビジネス – Cambria/Calibri、青アクセント",
	"preset_us-modern":   "米国モダン – Segoe UI、ミニマル",
	"preset_cn-official": "中国公文 – SimHei/SimSun（公文形式）、赤アクセント",
	"preset_cn-modern":   "中国モダン – Noto Sans SC、現代的",
	"preset_jp-formal":   "日本フォーマル – 游明朝/游ゴシック、ビジネス",
	"preset_eu-clean":    "欧州クリーン – Helvetica/Arial、ミニマル",
	"preset_kr-standard": "韓国標準 – Malgun Gothic/Nanum、クリーン",
	"preset_academic":    "アカデミック – Times New Roman、学術的",
	"preset_default":     "デフォルト – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "Markdown入力ファイルが見つかりません",
	"err_output_ext":           "出力パスは .docx で終わる必要があります",
	"err_style_not_found":      "スタイルテンプレートが見つかりません",
	"err_style_invalid_json":   "スタイルテンプレートが有効なJSONではありません",
	"err_style_missing_field":  "スタイルテンプレートに不足フィールドがあります",
	"err_style_bad_color":      "色は #RRGGBB 形式である必要があります",
	"err_style_bad_size":       "サイズは正の値である必要があります",
	"err_reading_file":         "ファイル読み込みエラー",
	"err_writing_file":         "ファイル書き込みエラー",
	"err_converting":           "変換に失敗しました",
	"err_tui":                  "TUIエラー",
}

// ---------------------------------------------------------------------------
// 한국어 (Korean)
// ---------------------------------------------------------------------------

var ko = map[string]string{
	"app_title":       "md2docx — Markdown → DOCX",
	"app_short":       "Markdown을 전문 DOCX 문서로 변환",

	"lang_select_title":   "언어 선택 / Select Language",
	"lang_select_help":    "↑↓ 이동   Enter 선택   Q 종료",
	"lang_saved":          "언어 저장됨: %s",
	"lang_saved_hint":     "변경: md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Markdown 입력 선택",
	"tui_output_title":    "md2docx — 출력 위치 선택",
	"tui_style_title":     "md2docx — 스타일 선택",
	"tui_style_custom":    "사용자 정의 JSON 파일...",
	"tui_style_default":   "기본값 사용",
	"tui_style_file_title":"md2docx — 스타일 템플릿 선택 (JSON)",
	"tui_confirm_title":   "md2docx — 변환 확인",
	"tui_confirm_input":   "입력",
	"tui_confirm_output":  "출력",
	"tui_confirm_style":   "스타일",
	"tui_confirm_mermaid": "Mermaid 다이어그램 렌더링",
	"tui_confirm_convert": "변환",
	"tui_confirm_back":    "뒤로",
	"tui_converting":      "변환 중...",
	"tui_done":            "완료",
	"tui_done_msg":        "완료: %s\n\n아무 키나 눌러 종료.",
	"tui_error_msg":       "오류: %v\n\nEsc로 종료.",
	"tui_filename_label":  "파일명",
	"tui_nav_help":        "↑↓ 이동 • Enter 선택 • Esc 뒤로",
	"tui_nav_quit":        "↑↓ 이동 • Enter 선택 • Esc 종료",
	"tui_nav_tab":         "↑↓ 이동 • Enter 확인 • Tab 전환 • Esc 뒤로",

	"cli_ok":             "성공",
	"cli_error":          "오류",
	"cli_bytes":          "바이트",
	"cli_template_created": "스타일 템플릿 생성됨: %s",
	"cli_no_agents":        "에이전트가 감지되지 않았습니다.",

	"preset_us-business": "미국 비즈니스 – Cambria/Calibri, 파란 강조",
	"preset_us-modern":   "미국 모던 – Segoe UI, 미니멀",
	"preset_cn-official": "중국 공문 – SimHei/SimSun (공문서), 빨간 강조",
	"preset_cn-modern":   "중국 모던 – Noto Sans SC",
	"preset_jp-formal":   "일본 포멀 – Yu Mincho/Yu Gothic, 비즈니스",
	"preset_eu-clean":    "유럽 클린 – Helvetica/Arial, 미니멀",
	"preset_kr-standard": "한국 표준 – Malgun Gothic/Nanum",
	"preset_academic":    "학술 – Times New Roman",
	"preset_default":     "기본 – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "Markdown 입력 파일을 찾을 수 없습니다",
	"err_output_ext":           "출력 경로는 .docx로 끝나야 합니다",
	"err_style_not_found":      "스타일 템플릿을 찾을 수 없습니다",
	"err_style_invalid_json":   "스타일 템플릿이 유효한 JSON이 아닙니다",
	"err_style_missing_field":  "스타일 템플릿에 필수 필드가 누락되었습니다",
	"err_style_bad_color":      "색상은 #RRGGBB 형식이어야 합니다",
	"err_style_bad_size":       "크기는 양수여야 합니다",
	"err_reading_file":         "파일 읽기 오류",
	"err_writing_file":         "파일 쓰기 오류",
	"err_converting":           "변환 실패",
	"err_tui":                  "TUI 오류",
}

// ---------------------------------------------------------------------------
// Español (Spanish)
// ---------------------------------------------------------------------------

var es = map[string]string{
	"app_title":       "md2docx — Markdown a DOCX",
	"app_short":       "Convierte Markdown en documentos DOCX profesionales",

	"lang_select_title":   "Seleccionar idioma / Select Language",
	"lang_select_help":    "↑↓ mover   Enter seleccionar   Q salir",
	"lang_saved":          "Idioma guardado: %s",
	"lang_saved_hint":     "Cambiar con: md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Seleccionar entrada Markdown",
	"tui_output_title":    "md2docx — Elegir salida",
	"tui_style_title":     "md2docx — Elegir estilo",
	"tui_style_custom":    "Archivo JSON personalizado...",
	"tui_style_default":   "Usar predeterminado",
	"tui_style_file_title":"md2docx — Seleccionar plantilla (JSON)",
	"tui_confirm_title":   "md2docx — Confirmar conversión",
	"tui_confirm_input":   "Entrada",
	"tui_confirm_output":  "Salida",
	"tui_confirm_style":   "Estilo",
	"tui_confirm_mermaid": "Renderizar diagramas Mermaid",
	"tui_confirm_convert": "Convertir",
	"tui_confirm_back":    "Volver",
	"tui_converting":      "Convirtiendo...",
	"tui_done":            "Hecho",
	"tui_done_msg":        "Hecho: %s\n\nPresione cualquier tecla para salir.",
	"tui_error_msg":       "Error: %v\n\nPresione Esc para salir.",
	"tui_filename_label":  "Nombre",
	"tui_nav_help":        "↑↓ mover • Enter seleccionar • Esc volver",
	"tui_nav_quit":        "↑↓ mover • Enter seleccionar • Esc salir",
	"tui_nav_tab":         "↑↓ mover • Enter confirmar • Tab cambiar • Esc volver",

	"cli_ok":             "OK",
	"cli_error":          "ERROR",
	"cli_bytes":          "bytes",
	"cli_template_created": "Plantilla creada: %s",
	"cli_no_agents":        "No se detectaron agentes.",

	"preset_us-business": "EEUU Negocios – Cambria/Calibri, acento azul",
	"preset_us-modern":   "EEUU Moderno – Segoe UI, tonos oscuros",
	"preset_cn-official": "China Oficial – SimHei/SimSun, acento rojo",
	"preset_cn-modern":   "China Moderno – Noto Sans SC",
	"preset_jp-formal":   "Japón Formal – Yu Mincho/Yu Gothic",
	"preset_eu-clean":    "Europa Limpio – Helvetica/Arial, minimalista",
	"preset_kr-standard": "Corea Estándar – Malgun Gothic/Nanum",
	"preset_academic":    "Académico – Times New Roman",
	"preset_default":     "Predeterminado – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "No se encontró el archivo Markdown",
	"err_output_ext":           "La salida debe terminar en .docx",
	"err_style_not_found":      "Plantilla de estilo no encontrada",
	"err_style_invalid_json":   "La plantilla no es JSON válido",
	"err_style_missing_field":  "Falta un campo en la plantilla",
	"err_style_bad_color":      "El color debe ser formato #RRGGBB",
	"err_style_bad_size":       "El tamaño debe ser positivo",
	"err_reading_file":         "Error al leer archivo",
	"err_writing_file":         "Error al escribir archivo",
	"err_converting":           "Falló la conversión",
	"err_tui":                  "Error de TUI",
}

// ---------------------------------------------------------------------------
// Português (Brasil)
// ---------------------------------------------------------------------------

var ptBR = map[string]string{
	"app_title":       "md2docx — Markdown para DOCX",
	"app_short":       "Converta Markdown em documentos DOCX profissionais",

	"lang_select_title":   "Selecionar idioma / Select Language",
	"lang_select_help":    "↑↓ mover   Enter selecionar   Q sair",
	"lang_saved":          "Idioma salvo: %s",
	"lang_saved_hint":     "Alterar com: md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Selecionar entrada Markdown",
	"tui_output_title":    "md2docx — Escolher saída",
	"tui_style_title":     "md2docx — Escolher estilo",
	"tui_style_custom":    "Arquivo JSON personalizado...",
	"tui_style_default":   "Usar padrão",
	"tui_style_file_title":"md2docx — Selecionar modelo (JSON)",
	"tui_confirm_title":   "md2docx — Confirmar conversão",
	"tui_confirm_input":   "Entrada",
	"tui_confirm_output":  "Saída",
	"tui_confirm_style":   "Estilo",
	"tui_confirm_mermaid": "Renderizar diagramas Mermaid",
	"tui_confirm_convert": "Converter",
	"tui_confirm_back":    "Voltar",
	"tui_converting":      "Convertendo...",
	"tui_done":            "Concluído",
	"tui_done_msg":        "Concluído: %s\n\nPressione qualquer tecla para sair.",
	"tui_error_msg":       "Erro: %v\n\nPressione Esc para sair.",
	"tui_filename_label":  "Nome",
	"tui_nav_help":        "↑↓ mover • Enter selecionar • Esc voltar",
	"tui_nav_quit":        "↑↓ mover • Enter selecionar • Esc sair",
	"tui_nav_tab":         "↑↓ mover • Enter confirmar • Tab alternar • Esc voltar",

	"cli_ok":             "OK",
	"cli_error":          "ERRO",
	"cli_bytes":          "bytes",
	"cli_template_created": "Modelo criado: %s",
	"cli_no_agents":        "Nenhum agente detectado.",

	"preset_us-business": "EUA Negócios – Cambria/Calibri, azul",
	"preset_us-modern":   "EUA Moderno – Segoe UI, minimal",
	"preset_cn-official": "China Oficial – SimHei/SimSun, vermelho",
	"preset_cn-modern":   "China Moderno – Noto Sans SC",
	"preset_jp-formal":   "Japão Formal – Yu Mincho/Yu Gothic",
	"preset_eu-clean":    "Europa Limpo – Helvetica/Arial, minimalista",
	"preset_kr-standard": "Coreia Padrão – Malgun Gothic/Nanum",
	"preset_academic":    "Acadêmico – Times New Roman",
	"preset_default":     "Padrão – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "Arquivo Markdown não encontrado",
	"err_output_ext":           "Saída deve terminar em .docx",
	"err_style_not_found":      "Modelo de estilo não encontrado",
	"err_style_invalid_json":   "Modelo não é JSON válido",
	"err_style_missing_field":  "Campo obrigatório ausente no modelo",
	"err_style_bad_color":      "Cor deve ser formato #RRGGBB",
	"err_style_bad_size":       "Tamanho deve ser positivo",
	"err_reading_file":         "Erro ao ler arquivo",
	"err_writing_file":         "Erro ao gravar arquivo",
	"err_converting":           "Falha na conversão",
	"err_tui":                  "Erro de TUI",
}

// ---------------------------------------------------------------------------
// Deutsch (German)
// ---------------------------------------------------------------------------

var de = map[string]string{
	"app_title":       "md2docx — Markdown zu DOCX",
	"app_short":       "Konvertiert Markdown in professionelle DOCX-Dokumente",

	"lang_select_title":   "Sprache wählen / Select Language",
	"lang_select_help":    "↑↓ navigieren   Enter auswählen   Q beenden",
	"lang_saved":          "Sprache gespeichert: %s",
	"lang_saved_hint":     "Ändern mit: md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Markdown-Eingabe wählen",
	"tui_output_title":    "md2docx — Ausgabe wählen",
	"tui_style_title":     "md2docx — Stil wählen",
	"tui_style_custom":    "Eigene JSON-Datei...",
	"tui_style_default":   "Standard verwenden",
	"tui_style_file_title":"md2docx — Stilvorlage wählen (JSON)",
	"tui_confirm_title":   "md2docx — Konvertierung bestätigen",
	"tui_confirm_input":   "Eingabe",
	"tui_confirm_output":  "Ausgabe",
	"tui_confirm_style":   "Stil",
	"tui_confirm_mermaid": "Mermaid-Diagramme rendern",
	"tui_confirm_convert": "Konvertieren",
	"tui_confirm_back":    "Zurück",
	"tui_converting":      "Konvertiere...",
	"tui_done":            "Fertig",
	"tui_done_msg":        "Fertig: %s\n\nBeliebige Taste zum Beenden.",
	"tui_error_msg":       "Fehler: %v\n\nEsc zum Beenden.",
	"tui_filename_label":  "Dateiname",
	"tui_nav_help":        "↑↓ navigieren • Enter auswählen • Esc zurück",
	"tui_nav_quit":        "↑↓ navigieren • Enter auswählen • Esc beenden",
	"tui_nav_tab":         "↑↓ navigieren • Enter bestätigen • Tab wechseln • Esc zurück",

	"cli_ok":             "OK",
	"cli_error":          "FEHLER",
	"cli_bytes":          "Bytes",
	"cli_template_created": "Stilvorlage erstellt: %s",
	"cli_no_agents":        "Keine Agenten erkannt.",

	"preset_us-business": "US Business – Cambria/Calibri, blau",
	"preset_us-modern":   "US Modern – Segoe UI, minimal",
	"preset_cn-official": "CN Offiziell – SimHei/SimSun, rot",
	"preset_cn-modern":   "CN Modern – Noto Sans SC",
	"preset_jp-formal":   "JP Formal – Yu Mincho/Yu Gothic",
	"preset_eu-clean":    "EU Clean – Helvetica/Arial, minimalistisch",
	"preset_kr-standard": "KR Standard – Malgun Gothic/Nanum",
	"preset_academic":    "Akademisch – Times New Roman",
	"preset_default":     "Standard – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "Markdown-Datei nicht gefunden",
	"err_output_ext":           "Ausgabepfad muss mit .docx enden",
	"err_style_not_found":      "Stilvorlage nicht gefunden",
	"err_style_invalid_json":   "Stilvorlage ist kein gültiges JSON",
	"err_style_missing_field":  "Pflichtfeld fehlt in Stilvorlage",
	"err_style_bad_color":      "Farbe muss im Format #RRGGBB sein",
	"err_style_bad_size":       "Größe muss positiv sein",
	"err_reading_file":         "Fehler beim Lesen der Datei",
	"err_writing_file":         "Fehler beim Schreiben der Datei",
	"err_converting":           "Konvertierung fehlgeschlagen",
	"err_tui":                  "TUI-Fehler",
}

// ---------------------------------------------------------------------------
// Français (French)
// ---------------------------------------------------------------------------

var fr = map[string]string{
	"app_title":       "md2docx — Markdown vers DOCX",
	"app_short":       "Convertit Markdown en documents DOCX professionnels",

	"lang_select_title":   "Choisir la langue / Select Language",
	"lang_select_help":    "↑↓ naviguer   Enter sélectionner   Q quitter",
	"lang_saved":          "Langue enregistrée : %s",
	"lang_saved_hint":     "Modifier avec : md2docx --lang <code>\n",

	"tui_input_title":     "md2docx — Sélectionner l'entrée Markdown",
	"tui_output_title":    "md2docx — Choisir la sortie",
	"tui_style_title":     "md2docx — Choisir le style",
	"tui_style_custom":    "Fichier JSON personnalisé...",
	"tui_style_default":   "Utiliser par défaut",
	"tui_style_file_title":"md2docx — Sélectionner le modèle (JSON)",
	"tui_confirm_title":   "md2docx — Confirmer la conversion",
	"tui_confirm_input":   "Entrée",
	"tui_confirm_output":  "Sortie",
	"tui_confirm_style":   "Style",
	"tui_confirm_mermaid": "Rendre les diagrammes Mermaid",
	"tui_confirm_convert": "Convertir",
	"tui_confirm_back":    "Retour",
	"tui_converting":      "Conversion...",
	"tui_done":            "Terminé",
	"tui_done_msg":        "Terminé : %s\n\nAppuyez sur une touche pour quitter.",
	"tui_error_msg":       "Erreur : %v\n\nAppuyez sur Échap pour quitter.",
	"tui_filename_label":  "Nom",
	"tui_nav_help":        "↑↓ naviguer • Enter sélectionner • Esc retour",
	"tui_nav_quit":        "↑↓ naviguer • Enter sélectionner • Esc quitter",
	"tui_nav_tab":         "↑↓ naviguer • Enter confirmer • Tab basculer • Esc retour",

	"cli_ok":             "OK",
	"cli_error":          "ERREUR",
	"cli_bytes":          "octets",
	"cli_template_created": "Modèle créé : %s",
	"cli_no_agents":        "Aucun agent détecté.",

	"preset_us-business": "US Business – Cambria/Calibri, accent bleu",
	"preset_us-modern":   "US Moderne – Segoe UI, tons sombres",
	"preset_cn-official": "CN Officiel – SimHei/SimSun, accent rouge",
	"preset_cn-modern":   "CN Moderne – Noto Sans SC",
	"preset_jp-formal":   "JP Formel – Yu Mincho/Yu Gothic",
	"preset_eu-clean":    "EU Épuré – Helvetica/Arial, minimaliste",
	"preset_kr-standard": "KR Standard – Malgun Gothic/Nanum",
	"preset_academic":    "Académique – Times New Roman",
	"preset_default":     "Défaut – Aptos Display/Cascadia Mono",

	"err_markdown_not_found":   "Fichier Markdown introuvable",
	"err_output_ext":           "La sortie doit se terminer par .docx",
	"err_style_not_found":      "Modèle de style introuvable",
	"err_style_invalid_json":   "Le modèle n'est pas un JSON valide",
	"err_style_missing_field":  "Champ obligatoire manquant dans le modèle",
	"err_style_bad_color":      "La couleur doit être au format #RRGGBB",
	"err_style_bad_size":       "La taille doit être positive",
	"err_reading_file":         "Erreur de lecture du fichier",
	"err_writing_file":         "Erreur d'écriture du fichier",
	"err_converting":           "Échec de la conversion",
	"err_tui":                  "Erreur TUI",
}
