package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Data structures for SRD content
type SRDEquipment struct {
	Name        string `json:"name"`
	Cost        string `json:"cost"`
	Weight      string `json:"weight"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Equippable  bool   `json:"equippable"`
	Slot        string `json:"slot"` // head, body, hands, feet, etc.
}

type SRDWeapon struct {
	Name        string `json:"name"`
	Cost        string `json:"cost"`
	Damage      string `json:"damage"`
	Weight      string `json:"weight"`
	Properties  string `json:"properties"`
	Description string `json:"description"`
}

type SRDSpell struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	School      string `json:"school"`
	CastingTime string `json:"casting_time"`
	Range       string `json:"range"`
	Components  string `json:"components"`
	Duration    string `json:"duration"`
	Description string `json:"description"`
	Classes     string `json:"classes"`
}

type SRDData struct {
	Equipment []SRDEquipment `json:"equipment"`
	Weapons   []SRDWeapon    `json:"weapons"`
	Spells    []SRDSpell     `json:"spells"`
}

// Character represents a D&D character
type Character struct {
	Name          string     `json:"name"`
	Race          string     `json:"race"`
	Class         string     `json:"class"`
	Level         int        `json:"level"`
	Abilities     Abilities  `json:"abilities"`
	Skills        []Skill    `json:"skills"`
	Equipment     []Item     `json:"equipment"`
	Weapons       []Weapon   `json:"weapons"`
	Spells        []Spell    `json:"spells"`
	Background    string     `json:"background"`
	Proficiencies []string   `json:"proficiencies"`
	Currency      Currency   `json:"currency"`
	Equipped      Equipped   `json:"equipped"`
}

type Equipped struct {
	Head    Item `json:"head"`
	Body    Item `json:"body"`
	Hands   Item `json:"hands"`
	Feet    Item `json:"feet"`
	Ring1   Item `json:"ring1"`
	Ring2   Item `json:"ring2"`
	Neck    Item `json:"neck"`
	MainHand Weapon `json:"mainHand"`
	OffHand Weapon `json:"offHand"`
}

type Currency struct {
	CP int `json:"cp"`
	SP int `json:"sp"`
	EP int `json:"ep"`
	GP int `json:"gp"`
	PP int `json:"pp"`
}

type Abilities struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type Skill struct {
	Name       string `json:"name"`
	Proficient bool   `json:"proficient"`
	Modifier   int    `json:"modifier"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Weight      string `json:"weight"`
	Cost        string `json:"cost"`
	Equipped    bool   `json:"equipped"`
	Slot        string `json:"slot"` // Where it can be equipped
}

type Weapon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Damage      string `json:"damage"`
	Properties  string `json:"properties"`
	Weight      string `json:"weight"`
	Cost        string `json:"cost"`
	Equipped    bool   `json:"equipped"`
	Hand        string `json:"hand"` // main or off
}

type Spell struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	School      string `json:"school"`
	Prepared    bool   `json:"prepared"`
}

// Model represents the application state
type Model struct {
	character     Character
	tabs          []string
	activeTab     int
	inputs        map[string]textinput.Model
	skillList     []Skill
	equipment     []Item
	weapons       []Weapon
	spells        []Spell
	proficiencies []string
	srdData       SRDData
	mode          string // "view" or "edit"
	message       string
	equipMode     string // "inventory" or "equipped"
	selectedItem  int    // Index of selected item in equipment list
	selectedWeapon int   // Index of selected weapon in weapons list
}

// Initialization
func initialModel() Model {
	char := Character{
		Name:   "New Character",
		Race:   "Human",
		Class:  "Fighter",
		Level:  1,
		Abilities: Abilities{
			Strength:     10,
			Dexterity:    10,
			Constitution: 10,
			Intelligence: 10,
			Wisdom:       10,
			Charisma:     10,
		},
		Skills:        []Skill{},
		Equipment:     []Item{},
		Weapons:       []Weapon{},
		Spells:        []Spell{},
		Background:    "Acolyte",
		Proficiencies: []string{},
		Currency: Currency{
			CP: 0,
			SP: 0,
			EP: 0,
			GP: 15, // Starting gold for most classes
			PP: 0,
		},
		Equipped: Equipped{},
	}

	// Define tabs
	tabs := []string{
		"Basic Info",
		"Abilities",
		"Skills",
		"Equipment",
		"Weapons",
		"Spells",
		"Background",
		"Proficiencies",
		"Currency",
	}

	// Initialize text inputs
	inputs := make(map[string]textinput.Model)
	inputs["name"] = createInput("Name: ", char.Name)
	inputs["race"] = createInput("Race: ", char.Race)
	inputs["class"] = createInput("Class: ", char.Class)
	inputs["level"] = createInput("Level: ", fmt.Sprintf("%d", char.Level))
	inputs["background"] = createInput("Background: ", char.Background)

	// Initialize ability inputs
	inputs["str"] = createInput("Strength: ", fmt.Sprintf("%d", char.Abilities.Strength))
	inputs["dex"] = createInput("Dexterity: ", fmt.Sprintf("%d", char.Abilities.Dexterity))
	inputs["con"] = createInput("Constitution: ", fmt.Sprintf("%d", char.Abilities.Constitution))
	inputs["int"] = createInput("Intelligence: ", fmt.Sprintf("%d", char.Abilities.Intelligence))
	inputs["wis"] = createInput("Wisdom: ", fmt.Sprintf("%d", char.Abilities.Wisdom))
	inputs["cha"] = createInput("Charisma: ", fmt.Sprintf("%d", char.Abilities.Charisma))

	// Initialize currency inputs
	inputs["cp"] = createInput("Copper: ", fmt.Sprintf("%d", char.Currency.CP))
	inputs["sp"] = createInput("Silver: ", fmt.Sprintf("%d", char.Currency.SP))
	inputs["ep"] = createInput("Electrum: ", fmt.Sprintf("%d", char.Currency.EP))
	inputs["gp"] = createInput("Gold: ", fmt.Sprintf("%d", char.Currency.GP))
	inputs["pp"] = createInput("Platinum: ", fmt.Sprintf("%d", char.Currency.PP))

	// Load SRD data
	srdData, err := loadSRDData()
	if err != nil {
		fmt.Printf("Error loading SRD data: %v\n", err)
	}

	return Model{
		character:     char,
		tabs:          tabs,
		activeTab:     0,
		inputs:        inputs,
		skillList:     char.Skills,
		equipment:     char.Equipment,
		weapons:       char.Weapons,
		spells:        char.Spells,
		proficiencies: char.Proficiencies,
		srdData:       srdData,
		mode:          "view",
		message:       "Welcome to D&D Character Editor!",
		equipMode:     "inventory",
		selectedItem:  -1,
		selectedWeapon: -1,
	}
}

func createInput(placeholder, value string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.SetValue(value)
	input.Focus()
	input.CharLimit = 50
	input.Width = 30
	return input
}

func loadSRDData() (SRDData, error) {
	var srdData SRDData
	
	// Load equipment data
	equipmentFile, err := os.ReadFile("srd_equipment.json")
	if err != nil {
		return srdData, fmt.Errorf("error reading equipment file: %v", err)
	}
	err = json.Unmarshal(equipmentFile, &srdData.Equipment)
	if err != nil {
		return srdData, fmt.Errorf("error parsing equipment JSON: %v", err)
	}
	
	// Load weapons data
	weaponsFile, err := os.ReadFile("srd_weapons.json")
	if err != nil {
		return srdData, fmt.Errorf("error reading weapons file: %v", err)
	}
	err = json.Unmarshal(weaponsFile, &srdData.Weapons)
	if err != nil {
		return srdData, fmt.Errorf("error parsing weapons JSON: %v", err)
	}
	
	// Load spells data
	spellsFile, err := os.ReadFile("srd_spells.json")
	if err != nil {
		return srdData, fmt.Errorf("error reading spells file: %v", err)
	}
	err = json.Unmarshal(spellsFile, &srdData.Spells)
	if err != nil {
		return srdData, fmt.Errorf("error parsing spells JSON: %v", err)
	}
	
	return srdData, nil
}

func saveCharacter(character Character) error {
	// Create characters directory if it doesn't exist
	err := os.MkdirAll("characters", 0755)
	if err != nil {
		return err
	}
	
	// Create filename from character name
	filename := fmt.Sprintf("characters/%s.json", strings.ReplaceAll(character.Name, " ", "_"))
	
	// Convert character to JSON
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	
	return nil
}

func loadCharacter(filename string) (Character, error) {
	var character Character
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return character, err
	}
	
	err = json.Unmarshal(data, &character)
	if err != nil {
		return character, err
	}
	
	return character, nil
}

func listSavedCharacters() ([]string, error) {
	var files []string
	
	err := filepath.Walk("characters", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	
	return files, nil
}

// UI Styles
var (
	tabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2)

	activeTabStyle = tabStyle.Copy().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A23C8"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A23C8")).
			Bold(true).
			Padding(0, 1)

	sectionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A23C8")).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2)

	selectedButtonStyle = buttonStyle.Copy().
			Background(lipgloss.Color("#5A23C8"))

	selectedItemStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	equippedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
)

// Init function for Bubble Tea
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update function for Bubble Tea
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			if m.activeTab == 3 && m.equipMode == "inventory" && len(m.equipment) > 0 {
				// In equipment tab, navigate items
				m.selectedItem = (m.selectedItem + 1) % len(m.equipment)
				return m, nil
			} else if m.activeTab == 4 && len(m.weapons) > 0 {
				// In weapons tab, navigate weapons
				m.selectedWeapon = (m.selectedWeapon + 1) % len(m.weapons)
				return m, nil
			} else {
				m.activeTab = (m.activeTab + 1) % len(m.tabs)
				m.selectedItem = -1
				m.selectedWeapon = -1
				return m, nil
			}
		case "left", "h", "p", "shift+tab":
			if m.activeTab == 3 && m.equipMode == "inventory" && len(m.equipment) > 0 {
				// In equipment tab, navigate items
				m.selectedItem = (m.selectedItem - 1 + len(m.equipment)) % len(m.equipment)
				return m, nil
			} else if m.activeTab == 4 && len(m.weapons) > 0 {
				// In weapons tab, navigate weapons
				m.selectedWeapon = (m.selectedWeapon - 1 + len(m.weapons)) % len(m.weapons)
				return m, nil
			} else {
				m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
				m.selectedItem = -1
				m.selectedWeapon = -1
				return m, nil
			}
		case "e":
			m.mode = "edit"
			return m, nil
		case "v":
			m.mode = "view"
			return m, nil
		case "s":
			if m.mode == "edit" {
				err := saveCharacter(m.character)
				if err != nil {
					m.message = fmt.Sprintf("Error saving character: %v", err)
				} else {
					m.message = "Character saved successfully!"
				}
			}
			return m, nil
		case "l":
			files, err := listSavedCharacters()
			if err != nil {
				m.message = fmt.Sprintf("Error listing characters: %v", err)
				return m, nil
			}
			if len(files) > 0 {
				// For simplicity, load the first character
				character, err := loadCharacter(files[0])
				if err != nil {
					m.message = fmt.Sprintf("Error loading character: %v", err)
				} else {
					m.character = character
					m.message = fmt.Sprintf("Loaded character: %s", character.Name)
					// Update inputs with loaded character data
					m.updateInputsFromCharacter()
				}
			} else {
				m.message = "No saved characters found"
			}
			return m, nil
		case "a":
			if m.mode == "edit" {
				if m.activeTab == 3 {
					// Add sample equipment
					newItem := Item{
						Name:        "Backpack",
						Description: "A backpack for carrying items",
						Quantity:    1,
						Weight:      "5 lb.",
						Cost:        "2 gp",
						Equipped:    false,
						Slot:        "back",
					}
					m.equipment = append(m.equipment, newItem)
					m.message = "Added backpack to equipment"
				} else if m.activeTab == 4 {
					// Add sample weapon
					newWeapon := Weapon{
						Name:        "Longsword",
						Description: "A versatile martial weapon",
						Damage:      "1d8 slashing",
						Properties:  "Versatile (1d10)",
						Weight:      "3 lb.",
						Cost:        "15 gp",
						Equipped:    false,
						Hand:        "",
					}
					m.weapons = append(m.weapons, newWeapon)
					m.message = "Added longsword to weapons"
				}
			}
			return m, nil
		case "enter":
			if m.mode == "edit" {
				return m.updateCharacterData()
			}
		case " ":
			if m.mode == "edit" {
				if m.activeTab == 3 && m.selectedItem >= 0 && m.selectedItem < len(m.equipment) {
					// Toggle equipment equipped status
					item := &m.equipment[m.selectedItem]
					item.Equipped = !item.Equipped
					
					if item.Equipped {
						// Equip the item to the appropriate slot
						switch item.Slot {
						case "head":
							m.character.Equipped.Head = *item
						case "body":
							m.character.Equipped.Body = *item
						case "hands":
							m.character.Equipped.Hands = *item
						case "feet":
							m.character.Equipped.Feet = *item
						case "ring":
							if m.character.Equipped.Ring1.Name == "" {
								m.character.Equipped.Ring1 = *item
							} else if m.character.Equipped.Ring2.Name == "" {
								m.character.Equipped.Ring2 = *item
							} else {
								// No ring slots available
								item.Equipped = false
								m.message = "No ring slots available"
							}
						case "neck":
							m.character.Equipped.Neck = *item
						}
					} else {
						// Unequip the item
						switch item.Slot {
						case "head":
							if m.character.Equipped.Head.Name == item.Name {
								m.character.Equipped.Head = Item{}
							}
						case "body":
							if m.character.Equipped.Body.Name == item.Name {
								m.character.Equipped.Body = Item{}
							}
						case "hands":
							if m.character.Equipped.Hands.Name == item.Name {
								m.character.Equipped.Hands = Item{}
							}
						case "feet":
							if m.character.Equipped.Feet.Name == item.Name {
								m.character.Equipped.Feet = Item{}
							}
						case "ring":
							if m.character.Equipped.Ring1.Name == item.Name {
								m.character.Equipped.Ring1 = Item{}
							} else if m.character.Equipped.Ring2.Name == item.Name {
								m.character.Equipped.Ring2 = Item{}
							}
						case "neck":
							if m.character.Equipped.Neck.Name == item.Name {
								m.character.Equipped.Neck = Item{}
							}
						}
					}
					
					m.message = fmt.Sprintf("%s %s", item.Name, map[bool]string{true: "equipped", false: "unequipped"}[item.Equipped])
				} else if m.activeTab == 4 && m.selectedWeapon >= 0 && m.selectedWeapon < len(m.weapons) {
					// Toggle weapon equipped status
					weapon := &m.weapons[m.selectedWeapon]
					weapon.Equipped = !weapon.Equipped
					
					if weapon.Equipped {
						// Equip the weapon
						if weapon.Hand == "main" || weapon.Hand == "" {
							if m.character.Equipped.MainHand.Name == "" {
								m.character.Equipped.MainHand = *weapon
								weapon.Hand = "main"
							} else if m.character.Equipped.OffHand.Name == "" {
								m.character.Equipped.OffHand = *weapon
								weapon.Hand = "off"
							} else {
								// No hand available
								weapon.Equipped = false
								m.message = "No hand available for weapon"
							}
						}
					} else {
						// Unequip the weapon
						if weapon.Hand == "main" && m.character.Equipped.MainHand.Name == weapon.Name {
							m.character.Equipped.MainHand = Weapon{}
						} else if weapon.Hand == "off" && m.character.Equipped.OffHand.Name == weapon.Name {
							m.character.Equipped.OffHand = Weapon{}
						}
						weapon.Hand = ""
					}
					
					m.message = fmt.Sprintf("%s %s", weapon.Name, map[bool]string{true: "equipped", false: "unequipped"}[weapon.Equipped])
				}
			}
			return m, nil
		case "i":
			if m.activeTab == 3 {
				m.equipMode = "inventory"
				m.message = "Viewing inventory"
			}
			return m, nil
		case "w":
			if m.activeTab == 3 {
				m.equipMode = "equipped"
				m.message = "Viewing equipped items"
			}
			return m, nil
		}
	}

	// Update the focused input field
	if cmd := m.updateInputs(msg); cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m *Model) updateInputsFromCharacter() {
	m.inputs["name"].SetValue(m.character.Name)
	m.inputs["race"].SetValue(m.character.Race)
	m.inputs["class"].SetValue(m.character.Class)
	m.inputs["level"].SetValue(fmt.Sprintf("%d", m.character.Level))
	m.inputs["background"].SetValue(m.character.Background)
	
	m.inputs["str"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Strength))
	m.inputs["dex"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Dexterity))
	m.inputs["con"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Constitution))
	m.inputs["int"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Intelligence))
	m.inputs["wis"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Wisdom))
	m.inputs["cha"].SetValue(fmt.Sprintf("%d", m.character.Abilities.Charisma))
	
	m.inputs["cp"].SetValue(fmt.Sprintf("%d", m.character.Currency.CP))
	m.inputs["sp"].SetValue(fmt.Sprintf("%d", m.character.Currency.SP))
	m.inputs["ep"].SetValue(fmt.Sprintf("%d", m.character.Currency.EP))
	m.inputs["gp"].SetValue(fmt.Sprintf("%d", m.character.Currency.GP))
	m.inputs["pp"].SetValue(fmt.Sprintf("%d", m.character.Currency.PP))
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	// Update basic info inputs
	if m.activeTab == 0 {
		for key := range map[string]bool{"name": true, "race": true, "class": true, "level": true} {
			var cmd tea.Cmd
			m.inputs[key], cmd = m.inputs[key].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update ability inputs
	if m.activeTab == 1 {
		for key := range map[string]bool{"str": true, "dex": true, "con": true, "int": true, "wis": true, "cha": true} {
			var cmd tea.Cmd
			m.inputs[key], cmd = m.inputs[key].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update background input
	if m.activeTab == 6 {
		var cmd tea.Cmd
		m.inputs["background"], cmd = m.inputs["background"].Update(msg)
		cmds = append(cmds, cmd)
	}
	
	// Update currency inputs
	if m.activeTab == 8 {
		for key := range map[string]bool{"cp": true, "sp": true, "ep": true, "gp": true, "pp": true} {
			var cmd tea.Cmd
			m.inputs[key], cmd = m.inputs[key].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

func (m Model) updateCharacterData() (tea.Model, tea.Cmd) {
	// Update basic info
	if m.activeTab == 0 {
		m.character.Name = m.inputs["name"].Value()
		m.character.Race = m.inputs["race"].Value()
		m.character.Class = m.inputs["class"].Value()
		if level, err := strconv.Atoi(m.inputs["level"].Value()); err == nil {
			m.character.Level = level
		}
	}

	// Update abilities
	if m.activeTab == 1 {
		if val, err := strconv.Atoi(m.inputs["str"].Value()); err == nil {
			m.character.Abilities.Strength = val
		}
		if val, err := strconv.Atoi(m.inputs["dex"].Value()); err == nil {
			m.character.Abilities.Dexterity = val
		}
		if val, err := strconv.Atoi(m.inputs["con"].Value()); err == nil {
			m.character.Abilities.Constitution = val
		}
		if val, err := strconv.Atoi(m.inputs["int"].Value()); err == nil {
			m.character.Abilities.Intelligence = val
		}
		if val, err := strconv.Atoi(m.inputs["wis"].Value()); err == nil {
			m.character.Abilities.Wisdom = val
		}
		if val, err := strconv.Atoi(m.inputs["cha"].Value()); err == nil {
			m.character.Abilities.Charisma = val
		}
	}

	// Update background
	if m.activeTab == 6 {
		m.character.Background = m.inputs["background"].Value()
	}
	
	// Update currency
	if m.activeTab == 8 {
		if val, err := strconv.Atoi(m.inputs["cp"].Value()); err == nil {
			m.character.Currency.CP = val
		}
		if val, err := strconv.Atoi(m.inputs["sp"].Value()); err == nil {
			m.character.Currency.SP = val
		}
		if val, err := strconv.Atoi(m.inputs["ep"].Value()); err == nil {
			m.character.Currency.EP = val
		}
		if val, err := strconv.Atoi(m.inputs["gp"].Value()); err == nil {
			m.character.Currency.GP = val
		}
		if val, err := strconv.Atoi(m.inputs["pp"].Value()); err == nil {
			m.character.Currency.PP = val
		}
	}

	// Update skills, equipment, weapons, spells, and proficiencies
	m.character.Skills = m.skillList
	m.character.Equipment = m.equipment
	m.character.Weapons = m.weapons
	m.character.Spells = m.spells
	m.character.Proficiencies = m.proficiencies

	m.message = "Character data updated!"
	return m, nil
}

// View function for Bubble Tea
func (m Model) View() string {
	// Render tabs
	var tabs []string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, activeTabStyle.Render(tab))
		} else {
			tabs = append(tabs, tabStyle.Render(tab))
		}
	}
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Render mode indicator
	modeIndicator := buttonStyle.Render("VIEW")
	if m.mode == "edit" {
		modeIndicator = selectedButtonStyle.Render("EDIT")
	}

	// Render active tab content
	var content string
	switch m.activeTab {
	case 0:
		content = m.renderBasicInfo()
	case 1:
		content = m.renderAbilities()
	case 2:
		content = m.renderSkills()
	case 3:
		content = m.renderEquipment()
	case 4:
		content = m.renderWeapons()
	case 5:
		content = m.renderSpells()
	case 6:
		content = m.renderBackground()
	case 7:
		content = m.renderProficiencies()
	case 8:
		content = m.renderCurrency()
	}

	// Render message
	message := messageStyle.Render(m.message)

	// Combine all components
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, tabsRow, "  ", modeIndicator),
		"",
		content,
		"",
		message,
		"",
		"Press ←/→ to switch tabs, e to edit, v to view, s to save, l to load, q to quit",
	)
}

func (m Model) renderBasicInfo() string {
	basicInfo := titleStyle.Render("Character Basics") + "\n\n"
	basicInfo += m.inputs["name"].View() + "\n"
	basicInfo += m.inputs["race"].View() + "\n"
	basicInfo += m.inputs["class"].View() + "\n"
	basicInfo += m.inputs["level"].View() + "\n"
	return sectionStyle.Render(basicInfo)
}

func (m Model) renderAbilities() string {
	abilities := titleStyle.Render("Ability Scores") + "\n\n"
	abilities += m.inputs["str"].View() + "  "
	abilities += m.inputs["dex"].View() + "\n"
	abilities += m.inputs["con"].View() + "  "
	abilities += m.inputs["int"].View() + "\n"
	abilities += m.inputs["wis"].View() + "  "
	abilities += m.inputs["cha"].View() + "\n"
	return sectionStyle.Render(abilities)
}

func (m Model) renderSkills() string {
	skills := titleStyle.Render("Skills") + "\n\n"
	for _, skill := range m.skillList {
		proficiency := " "
		if skill.Proficient {
			proficiency = "✓"
		}
		skills += fmt.Sprintf("[%s] %s: %+d\n", proficiency, skill.Name, skill.Modifier)
	}
	skills += "\nPress 'a' to add a new skill"
	return sectionStyle.Render(skills)
}

func (m Model) renderEquipment() string {
	equipmentView := titleStyle.Render("Equipment") + "\n\n"
	
	// Mode selector
	inventoryBtn := buttonStyle.Render("Inventory")
	equippedBtn := buttonStyle.Render("Equipped")
	if m.equipMode == "inventory" {
		inventoryBtn = selectedButtonStyle.Render("Inventory")
	} else {
		equippedBtn = selectedButtonStyle.Render("Equipped")
	}
	equipmentView += lipgloss.JoinHorizontal(lipgloss.Top, inventoryBtn, " ", equippedBtn) + "\n\n"
	
	if m.equipMode == "inventory" {
		// Show inventory
		if len(m.equipment) == 0 {
			equipmentView += "No equipment in inventory\n"
		} else {
			for i, item := range m.equipment {
				itemStr := fmt.Sprintf("%d. %s (%s, %s)", i+1, item.Name, item.Cost, item.Weight)
				if item.Equipped {
					itemStr = equippedStyle.Render(itemStr + " [EQUIPPED]")
				}
				
				if i == m.selectedItem {
					itemStr = selectedItemStyle.Render(itemStr)
				}
				
				equipmentView += itemStr + "\n"
			}
		}
		equipmentView += "\nPress Space to equip/unequip, a to add item, i/w to switch view"
	} else {
		// Show equipped items
		equipmentView += "Head: " + m.character.Equipped.Head.Name + "\n"
		equipmentView += "Body: " + m.character.Equipped.Body.Name + "\n"
		equipmentView += "Hands: " + m.character.Equipped.Hands.Name + "\n"
		equipmentView += "Feet: " + m.character.Equipped.Feet.Name + "\n"
		equipmentView += "Ring 1: " + m.character.Equipped.Ring1.Name + "\n"
		equipmentView += "Ring 2: " + m.character.Equipped.Ring2.Name + "\n"
		equipmentView += "Neck: " + m.character.Equipped.Neck.Name + "\n"
		equipmentView += "\nPress i/w to switch view"
	}
	
	return sectionStyle.Render(equipmentView)
}

func (m Model) renderWeapons() string {
	weapons := titleStyle.Render("Weapons") + "\n\n"
	
	if len(m.weapons) == 0 {
		weapons += "No weapons\n"
	} else {
		for i, weapon := range m.weapons {
			weaponStr := fmt.Sprintf("%d. %s (%s, %s)", i+1, weapon.Name, weapon.Damage, weapon.Cost)
			if weapon.Equipped {
				weaponStr = equippedStyle.Render(weaponStr + " [EQUIPPED]")
			}
			
			if i == m.selectedWeapon {
				weaponStr = selectedItemStyle.Render(weaponStr)
			}
			
			weapons += weaponStr + "\n"
		}
	}
	weapons += "\nPress Space to equip/unequip, a to add weapon"
	return sectionStyle.Render(weapons)
}

func (m Model) renderSpells() string {
	spells := titleStyle.Render("Spells") + "\n\n"
	for _, spell := range m.spells {
		prepared := " "
		if spell.Prepared {
			prepared = "✓"
		}
		spells += fmt.Sprintf("[%s] Level %d: %s (%s)\n", prepared, spell.Level, spell.Name, spell.School)
	}
	spells += "\nPress 'a' to add a new spell from SRD"
	return sectionStyle.Render(spells)
}

func (m Model) renderBackground() string {
	bg := titleStyle.Render("Background & Traits") + "\n\n"
	bg += m.inputs["background"].View() + "\n"
	return sectionStyle.Render(bg)
}

func (m Model) renderProficiencies() string {
	profs := titleStyle.Render("Proficiencies & Languages") + "\n\n"
	for i, prof := range m.proficiencies {
		profs += fmt.Sprintf("%d. %s\n", i+1, prof)
	}
	profs += "\nPress 'a' to add a new proficiency"
	return sectionStyle.Render(profs)
}

func (m Model) renderCurrency() string {
	currency := titleStyle.Render("Currency") + "\n\n"
	currency += m.inputs["cp"].View() + "  "
	currency += m.inputs["sp"].View() + "\n"
	currency += m.inputs["ep"].View() + "  "
	currency += m.inputs["gp"].View() + "\n"
	currency += m.inputs["pp"].View() + "\n"
	return sectionStyle.Render(currency)
}

// Main function
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}

