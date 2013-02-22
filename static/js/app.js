$().ready(function() {
    $('#filter_show_none').click(function(e) {
        e.preventDefault();
        $('.filter input[type="checkbox"]').removeAttr('checked');
        $('.filter input[type="checkbox"]').trigger('change');
    });
    $('#filter_show_all').click(function(e) {
        e.preventDefault();
        $('.filter input[type="checkbox"]').attr('checked', 'checked');
        $('.filter input[type="checkbox"]').trigger('change');
    })

    $('.filter input').change(function() {
        swap_org($(this).attr('value'));
    })


})

function swap_org(org) {
    $('div.project[data-org="' + org + '"]').toggle();
}