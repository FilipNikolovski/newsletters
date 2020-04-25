import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import qs from "qs";
import { Box, Button, Select, FormField, TextInput } from "grommet";

import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import ButtonWithLoader from "../ui/ButtonWithLoader";

const campaignValidation = object().shape({
  name: string()
    .required()
    .max(191, "The name must not exceed 191 characters."),
  template_name: string()
    .required()
    .max(191, "The template name must not exceed 191 characters."),
});

const CreateForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState("");
  const [options, setOptions] = useState({
    collection: [],
    url: "/api/templates",
  });
  const callApi = async () => {
    const res = await axios(options.url);
    let url = "";
    if (res.data.next_token) {
      url = `/api/templates?next_token=${encodeURIComponent(
        res.data.next_token
      )}`;
    }

    setOptions({
      collection: [...options.collection, ...res.data.collection],
      url: url,
    });
  };

  useEffect(() => {
    callApi();
  }, []);

  const onMore = () => {
    if (options.url) {
      callApi();
    }
  };

  const onChange = ({ value: nextSelected }) => {
    setFieldValue("template_name", nextSelected.name);
    setSelected(nextSelected);
  };

  return (
    <Box
      direction="column"
      fill
      margin={{ left: "medium", right: "medium", bottom: "medium" }}
    >
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="name" label="Name">
            <TextInput
              name="name"
              onChange={handleChange}
              placeholder="Campaign name"
            />
            <ErrorMessage name="name" />
          </FormField>
          <FormField htmlFor="template_name" label="Choose template">
            <Select
              placeholder="select a template..."
              value={selected}
              labelKey="name"
              valueKey="name"
              options={options.collection}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />

            <ErrorMessage name="template_name" />
          </FormField>

          <Box direction="row" alignSelf="end" margin={{ top: "large" }}>
            <Box margin={{ right: "small" }}>
              <Button label="Cancel" onClick={() => hideModal()} />
            </Box>
            <Box>
              <ButtonWithLoader
                type="submit"
                primary
                disabled={isSubmitting}
                label="Save Campaign"
              />
            </Box>
          </Box>
        </Box>
      </form>
    </Box>
  );
};

CreateForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  setFieldValue: PropTypes.func,
};

const CreateCampaign = ({ callApi, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        let data = {
          name: values.name,
          template_name: values.template_name,
        };

        await axios.post("/api/campaigns", qs.stringify(data));

        createNotification("Campaign has been created successfully.");

        //done submitting, set submitting to false
        setSubmitting(false);
        await callApi({ url: "/api/campaigns" });

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to create campaign. Please try again.";

          createNotification(msg, "status-error");

          //done submitting, set submitting to false
          setSubmitting(false);
        }
      }
    };

    await postForm();

    return;
  };

  return (
    <Box direction="row">
      <Formik
        initialValues={{ name: "", template_name: "" }}
        onSubmit={handleSubmit}
        validationSchema={campaignValidation}
      >
        {(props) => <CreateForm {...props} hideModal={hideModal} />}
      </Formik>
    </Box>
  );
};

CreateCampaign.propTypes = {
  callApi: PropTypes.func,
  hideModal: PropTypes.func,
};

export default CreateCampaign;
