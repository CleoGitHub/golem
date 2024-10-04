package domainbuilder

// func (b *domainBuilder) buildModel(ctx context.Context, modelDefinition *coredomaindefinition.Model) *domainBuilder {
// 	if b.Err != nil {
// 		return b
// 	}

// 	m := &model.Model{
// 		JsonName: modelDefinition.Name,
// 		Struct: &model.Struct{
// 			Name:   GetModelName(ctx, modelDefinition),
// 			Fields: []*model.Field{},
// 		},
// 		Activable: modelDefinition.Activable,
// 	}

// 	b.ModelUsecaseStruct[m] = &model.Struct{}

// 	modelFieldNames := []string{}
// 	for _, f := range b.DefaultModelFields {
// 		modelFieldNames = append(modelFieldNames, f.Name)
// 		field, err := b.FieldDefinitionToField(ctx, f)
// 		if err != nil {
// 			b.Err = merror.Stack(err)
// 			return b
// 		}
// 		m.Struct.Fields = append(m.Struct.Fields, field)
// 	}

// 	// Add activable field if model is activable
// 	if modelDefinition.Activable {
// 		field, err := b.FieldDefinitionToField(ctx, &coredomaindefinition.Field{
// 			Name: "active",
// 			Type: coredomaindefinition.PrimitiveTypeBool,
// 		})
// 		if err != nil {
// 			b.Err = merror.Stack(err)
// 			return b
// 		}
// 		m.Struct.Fields = append(m.Struct.Fields, field)

// 		field = field.Copy()
// 		b.ModelUsecaseStruct[m].Fields = append(b.ModelUsecaseStruct[m].Fields, field)
// 	}

// 	// Add default fields to definition
// 	for _, field := range modelDefinition.Fields {
// 		if slices.Contains(modelFieldNames, field.Name) {
// 			b.Err = merror.Stack(NewErrDefaultFiedlRedefined(field.Name))
// 			return b
// 		}
// 		f, err := b.FieldDefinitionToField(ctx, field)
// 		if err != nil {
// 			b.Err = merror.Stack(err)
// 			return b
// 		}
// 		m.Struct.Fields = append(m.Struct.Fields, f)
// 		b.FieldToValidationRules[f] = field.Validations

// 		f = f.Copy()
// 		validationsValues := []string{}

// 		for _, validation := range field.Validations {
// 			switch validation.Rule {
// 			case coredomaindefinition.ValidationRuleRequired:
// 				validationsValues = append(validationsValues, "required")
// 			case coredomaindefinition.ValidationRuleEmail:
// 				validationsValues = append(validationsValues, "email")
// 			case coredomaindefinition.ValidationRuleUUID:
// 				validationsValues = append(validationsValues, "uuid")
// 			case coredomaindefinition.ValidationRuleHexColor:
// 				validationsValues = append(validationsValues, "hexcolor")
// 			case coredomaindefinition.ValidationRuleGT:
// 				if s, ok := validation.Value.(string); ok {
// 					validationsValues = append(validationsValues, "gt:"+s)
// 				} else {
// 					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
// 				}
// 			case coredomaindefinition.ValidationRuleGTE:
// 				if s, ok := validation.Value.(string); ok {
// 					validationsValues = append(validationsValues, "gte:"+s)
// 				} else {
// 					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
// 				}
// 			case coredomaindefinition.ValidationRuleLT:
// 				if s, ok := validation.Value.(string); ok {
// 					validationsValues = append(validationsValues, "lt:"+s)
// 				} else {
// 					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
// 				}
// 			case coredomaindefinition.ValidationRuleLTE:
// 				if s, ok := validation.Value.(string); ok {
// 					validationsValues = append(validationsValues, "lte:"+s)
// 				} else {
// 					b.Err = NewErrValidationValueExpectedType(string(validation.Rule), "string")
// 				}
// 			}
// 		}

// 		if len(validationsValues) > 0 {
// 			f.Tags = append(f.Tags, &model.Tag{
// 				Name:   "validate",
// 				Values: validationsValues,
// 			})
// 		}

// 		b.ModelUsecaseStruct[m].Fields = append(b.ModelUsecaseStruct[m].Fields, f)
// 	}

// 	b.ModelDefinitionToModel[modelDefinition] = m
// 	b.ModelToModelDefinition[m] = modelDefinition
// 	b.Domain.Models = append(b.Domain.Models, m)

// 	return b
// }
